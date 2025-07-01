package nats

import (
	"context"
	"log"
	"time"

	vaultStreamConfig "github.com/jurshsmith/vaultstream/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connect establishes a connection to the NATS server, ensures the stream is created/updated,
// and returns a JetStream API (new API) along with the underlying NATS connection.
func Connect() (jetstream.JetStream, *nats.Conn) {
	natsURL := vaultStreamConfig.VaultStreamNatsURL()
	natsPassword := vaultStreamConfig.VaultStreamNatsPassword()

	opts := []nats.Option{nats.UserInfo("admin", natsPassword)}

	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Fatalf("Failed connecting to NATS: %v", err)
	}

	if _, err := createOrUpdateStream(nc); err != nil {
		log.Fatalf("Error creating/updating stream: %v", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Error creating new JetStream API: %v", err)
	}

	return js, nc
}

// createOrUpdateStream ensures a stream named "vaultstream-streams" exists with the given subjects.
func createOrUpdateStream(nc *nats.Conn) (*nats.StreamInfo, error) {
	// Create a context with a timeout for the stream operation.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jsc, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Monolith stream config for all VaultStream Signer streams
	streamConfig := &nats.StreamConfig{
		Name:      vaultStreamConfig.EventsStreamName(),
		Subjects:  []string{"records.>", "keys.>"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxMsgs:   -1, // No limit on the number of messages.
		MaxBytes:  -1, // No limit on the total size of messages.
		MaxAge:    0,  // Retain messages indefinitely.
		Replicas:  1,
	}

	// Attempt to add the stream. If it already exists, update its configuration.
	streamInfo, err := jsc.AddStream(streamConfig, nats.Context(ctx))
	if err != nil {
		streamInfo, err = jsc.UpdateStream(streamConfig, nats.Context(ctx))
		if err != nil {
			return nil, err
		}
	}

	return streamInfo, nil
}
