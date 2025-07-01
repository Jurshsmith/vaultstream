package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/jurshsmith/vaultstream/config"
	"github.com/jurshsmith/vaultstream/database"
	"github.com/jurshsmith/vaultstream/logger"
	"github.com/jurshsmith/vaultstream/nats"
	"github.com/jurshsmith/vaultstream/types"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

var log *zap.Logger

func main() {
	log = logger.New()
	defer log.Sync()

	log.Info("Records Service running...")

	config.Setup()

	batchSize := config.RecordsBatchSize()

	dbClient := database.Connect()
	defer dbClient.Close()
	log.Debug("Database connection established")

	jetstreamClient, natsConn := nats.Connect()
	defer natsConn.Close()
	log.Debug("NATS JetStream connection established")

	mainContext := context.Background()

	var waitGroup sync.WaitGroup
	semaphoreQueue := make(chan struct{}, config.RecordsMaxConcurrency())

	totalBatches := config.TotalRecordBatches()

	for batchID := 1; batchID <= totalBatches; batchID++ {
		semaphoreQueue <- struct{}{} // acquire
		waitGroup.Add(1)

		go func(batchID int) {
			defer waitGroup.Done()
			defer func() { <-semaphoreQueue }() // release

			dbRecords, err := dbClient.Record.
				Query().
				Where(func(s *sql.Selector) {
					s.Where(sql.ExprP("mod((id-1), $1)+1 = $2", totalBatches, batchID))
				}).
				All(mainContext)

			if err != nil {
				log.Error("Batch query error", zap.Int("batchID", batchID), zap.Error(err))
				return
			}

			records := dbRecordsToBytes(dbRecords)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			subject := fmt.Sprintf("records.%d", batchID)
			pubAck, err := jetstreamClient.Publish(ctx, subject, records, jetstream.WithMsgID(subject))
			if err != nil {
				log.Fatal("Error publishing message", zap.String("subject", subject), zap.Error(err))
			}

			log.Debug("Published message", zap.String("stream", pubAck.Stream), zap.Uint64("sequence", pubAck.Sequence))
			log.Debug("Batch published", zap.Int("batchID", batchID), zap.Int("recordCount", len(dbRecords)))
		}(batchID)
	}

	waitGroup.Wait()
	log.Info("All records enqueued successfully!", zap.Int("batchSize", batchSize))
}

func dbRecordsToBytes(records []*database.Record) []byte {
	var recordList []types.Record

	for _, r := range records {
		recordList = append(recordList, types.Record{
			ID:         r.ID,
			InsertedAt: r.InsertedAt,
		})
	}

	data, err := json.Marshal(recordList)
	if err != nil {
		log.Error("Error marshaling records", zap.Error(err))
		return []byte{}
	}

	return data
}
