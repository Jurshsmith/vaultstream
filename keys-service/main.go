package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jurshsmith/vaultstream/config"
	"github.com/jurshsmith/vaultstream/logger"
	vaultStreamNats "github.com/jurshsmith/vaultstream/nats"
	"github.com/jurshsmith/vaultstream/types"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

var log *zap.Logger

func main() {
	log = logger.New()
	defer log.Sync()

	log.Info("Keys Service running...")

	config.Setup()

	totalKeys := config.TotalKeys()

	jetstreamClient, natsConn := vaultStreamNats.Connect()
	defer natsConn.Close()

	keys, err := generateKeys(totalKeys)
	if err != nil {
		log.Fatal("Error generating keys", zap.Error(err))
	}

	enqueueAllKeys(jetstreamClient, keys)
}

func enqueueAllKeys(jetstreamClient jetstream.JetStream, keys []*types.Key) {
	var waitGroup sync.WaitGroup
	semaphoreQueue := make(chan struct{}, config.KeysMaxConcurrency())

	for keyID := range keys {
		semaphoreQueue <- struct{}{} // acquire
		waitGroup.Add(1)

		go func(keyID int) {
			defer waitGroup.Done()
			defer func() { <-semaphoreQueue }() // release

			key := keys[keyID]

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			enqueueKey(jetstreamClient, key, ctx)

		}(keyID)
	}

	waitGroup.Wait()

	log.Info("All keys initially enqueued successfully!", zap.Int("totalKeys", len(keys)))
}

func enqueueKey(jetstreamClient jetstream.JetStream, key *types.Key, ctx context.Context) {
	keyInBytes, err := json.Marshal(key)
	if err != nil {
		log.Fatal("Error marshaling key", zap.Error(err))
	}

	subject := fmt.Sprintf("keys.%d", key.ID)
	pubAck, err := jetstreamClient.Publish(ctx, subject, keyInBytes, jetstream.WithMsgID(subject))
	if err != nil {
		log.Fatal("Error publishing message", zap.String("subject", subject), zap.Error(err))
	}

	log.Debug("Published message", zap.String("stream", pubAck.Stream), zap.Uint64("sequence", pubAck.Sequence))
}

func generateKeys(totalKeys int) ([]*types.Key, error) {
	keys := make([]*types.Key, totalKeys)
	for i := range totalKeys {
		key, err := generateKey(i + 1)
		if err != nil {
			return nil, fmt.Errorf("failed generating key %d: %w", i+1, err)
		}
		keys[i] = key
	}
	return keys, nil
}

func generateKey(id int) (*types.Key, error) {
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed generating key %d: %w", id, err)
	}
	derBytes, err := x509.MarshalECPrivateKey(ecdsaKey)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling key %d: %w", id, err)
	}
	encodedKey := base64.StdEncoding.EncodeToString(derBytes)
	k := &types.Key{
		ID:         id,
		Value:      encodedKey,
		IsInUse:    false,
		LastUsedAt: time.Unix(0, 0), // initial value indicating never used.
	}
	return k, nil
}
