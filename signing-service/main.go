package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jurshsmith/vaultstream/config"
	"github.com/jurshsmith/vaultstream/database"
	"github.com/jurshsmith/vaultstream/logger"
	"github.com/jurshsmith/vaultstream/nats"
	"github.com/jurshsmith/vaultstream/types"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var log *zap.Logger

func main() {
	log = logger.New()
	defer log.Sync()

	log.Info("Signing Service running...")
	startTime := time.Now()

	// Load configuration
	config.Setup()

	// Setup database and NATS JetStream connections
	dbClient := database.Connect()
	defer dbClient.Close()
	log.Debug("Database connection established")

	ctx := context.Background()

	jetstreamClient, natsConn := nats.Connect()
	defer natsConn.Close()
	log.Debug("NATS JetStream connection established")

	recordsConsumerName := "signing-records-consumer"
	keysConsumerName := "signing-keys-consumer"

	// Create (or update) consumer for the records stream.
	recordsConsumerConfig := &jetstream.ConsumerConfig{
		Durable:       recordsConsumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "records.>",
		DeliverPolicy: jetstream.DeliverAllPolicy,
	}
	recordsConsumer, err := jetstreamClient.CreateOrUpdateConsumer(ctx, config.EventsStreamName(), *recordsConsumerConfig)
	if err != nil {
		log.Fatal("Error creating/updating records consumer", zap.Error(err))
	}

	// Create (or update) consumer for the keys stream.
	keysConsumerConfig := &jetstream.ConsumerConfig{
		Durable:       keysConsumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "keys.>",
		DeliverPolicy: jetstream.DeliverAllPolicy,
	}
	keysConsumer, err := jetstreamClient.CreateOrUpdateConsumer(ctx, config.EventsStreamName(), *keysConsumerConfig)
	if err != nil {
		log.Fatal("Error creating/updating keys consumer", zap.Error(err))
	}

	var (
		totalSignedRecords   int64 // use atomic counter
		expectedTotalBatches = config.TotalRecordBatches()
		batchesEnqueuedSoFar = 0
		procWG               sync.WaitGroup
	)

	// Main loop: keep pulling messages until we've signed enough records.
	for {
		if batchesEnqueuedSoFar == expectedTotalBatches {
			break
		}

		// Pull one records message (blocking until available)
		recordsMsg, err := recordsConsumer.Next()
		if err != nil {
			log.Error("Error fetching records message", zap.Error(err))
			continue
		}

		// Decode the batch of records.
		var records []types.Record
		if err := json.Unmarshal(recordsMsg.Data(), &records); err != nil {
			log.Error("Error unmarshaling records", zap.Error(err))
			recordsMsg.Nak() // Negative acknowledgment so it can be retried
			continue
		}
		log.Debug("Fetched records", zap.Int("RecordsBatchSize", len(records)))

		// Pull one free key message (blocking until available)
		keyMsg, err := keysConsumer.Next()
		if err != nil {
			log.Error("Error fetching key message", zap.Error(err))
			recordsMsg.Nak()
			continue
		}

		var key types.Key
		if err := json.Unmarshal(keyMsg.Data(), &key); err != nil {
			log.Error("Error unmarshaling key", zap.Error(err))
			keyMsg.Nak()
			recordsMsg.Nak()
			continue
		}
		log.Debug("Fetched free key", zap.Int("keyID", key.ID))

		// Spawn a goroutine to process signing and bulk insertion.
		procWG.Add(1)
		go func(recordsMsg jetstream.Msg, keyMsg jetstream.Msg, records []types.Record, key types.Key) {
			defer procWG.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Sign each record concurrently using the key.
			signatures, err := signRecords(records, &key)
			if err != nil {
				log.Error("Error signing records", zap.Error(err))
				return // Do not ack; message will be re-delivered.
			}

			// Bulk insert signatures into the database.
			if err := insertSignatures(ctx, dbClient, signatures); err != nil {
				log.Error("Error inserting signatures", zap.Error(err))
				return
			}

			// After successful DB insert, acknowledge the key message and re-publish the key.
			if err := keyMsg.Nak(); err != nil {
				log.Error("Error re-enqueueing key", zap.Error(err))
			}

			if err := recordsMsg.Ack(); err != nil {
				log.Error("Error acknowledging records being signed", zap.Error(err))
			}

			// Update the counter.
			atomic.AddInt64(&totalSignedRecords, int64(len(records)))
			log.Info("Batch processed", zap.Int64("totalRecordsSigned", atomic.LoadInt64(&totalSignedRecords)))
		}(recordsMsg, keyMsg, records, key)

		batchesEnqueuedSoFar++
	}

	// Wait for all processing goroutines to complete.
	procWG.Wait()

	elapsedTime := time.Since(startTime)
	log.Info("Signing service completed signing all records", zap.Int64("totalSigned", atomic.LoadInt64(&totalSignedRecords)), zap.Duration("elapsed", elapsedTime))
}

// signRecords spawns goroutines to sign each record concurrently using the provided key.
// This version pre-allocates a slice and assigns each signature by its index,
// preserving the input order.
func signRecords(records []types.Record, key *types.Key) ([]types.Signature, error) {
	sigs := make([]types.Signature, len(records))
	var wg sync.WaitGroup
	for i, rec := range records {
		wg.Add(1)
		go func(i int, r types.Record) {
			defer wg.Done()
			sigs[i] = signRecord(r, key)
		}(i, rec)
	}
	wg.Wait()
	return sigs, nil
}

// signRecord simulates the creation of a signature by hashing the record ID and key value.
func signRecord(record types.Record, key *types.Key) types.Signature {
	data := fmt.Sprintf("%d:%s", record.ID, key.Value)
	hash := sha256.Sum256([]byte(data))
	signatureValue := base64.StdEncoding.EncodeToString(hash[:])
	return types.Signature{
		RecordID: record.ID,
		KeyID:    key.ID,
		Value:    signatureValue,
	}
}

// insertSignatures performs a bulk insert of signatures using the ent ORM client.
func insertSignatures(ctx context.Context, client *database.Client, sigs []types.Signature) error {
	log.Info("Inserting batch of signatures into the DB", zap.Int("InsertedSignaturesBatchSize", len(sigs)))

	const (
		chunkSize   = 10000 // maximum number of rows per bulk insert
		maxParallel = 2     // maximum number of concurrent bulk inserts
	)

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(maxParallel)

	n := len(sigs)
	numChunks := (n + chunkSize - 1) / chunkSize

	for i := range numChunks {
		start := i * chunkSize
		end := min(start+chunkSize, n)
		chunk := sigs[start:end]

		// Capture chunk in a local variable for the goroutine.
		chunkCopy := chunk
		eg.Go(func() error {
			bulk := make([]*database.SignatureCreate, 0, len(chunkCopy))
			for _, sig := range chunkCopy {
				bulk = append(bulk, client.Signature.Create().
					SetRecordID(sig.RecordID).
					SetKeyID(sig.KeyID).
					SetValue(sig.Value))
			}
			if _, err := client.Signature.CreateBulk(bulk...).Save(ctx); err != nil {
				return fmt.Errorf("failed to bulk insert signatures: %w", err)
			}
			return nil
		})
	}

	return eg.Wait()
}
