// main_test.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jurshsmith/vaultstream/database"
	"github.com/jurshsmith/vaultstream/types"
	"github.com/nats-io/nats.go/jetstream"
)

// --- Tests for dbRecordsToBytes ---

func TestDBRecordsToBytes(t *testing.T) {
	// Use a fixed time to avoid issues with sub-second differences.
	now := time.Now().UTC().Truncate(time.Second)
	records := []*database.Record{
		{ID: 1, InsertedAt: now},
		{ID: 2, InsertedAt: now.Add(time.Second)},
	}
	data := dbRecordsToBytes(records)

	var recordsOut []types.Record
	err := json.Unmarshal(data, &recordsOut)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(recordsOut) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(recordsOut))
	}
	if recordsOut[0].ID != 1 || !recordsOut[0].InsertedAt.Equal(now) {
		t.Errorf("Unexpected first record: %+v", recordsOut[0])
	}
	if recordsOut[1].ID != 2 || !recordsOut[1].InsertedAt.Equal(now.Add(time.Second)) {
		t.Errorf("Unexpected second record: %+v", recordsOut[1])
	}
}

func TestDBRecordsToBytesEmpty(t *testing.T) {
	data := dbRecordsToBytes([]*database.Record{})
	var recordsOut []types.Record
	err := json.Unmarshal(data, &recordsOut)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if len(recordsOut) != 0 {
		t.Errorf("Expected 0 records, got %d", len(recordsOut))
	}
}

// --- Refactored Batch Processing Logic and Its Test ---
//
// To test the batch processing that occurs concurrently in main,
// we recommend refactoring the logic in the goroutine into a separate function.
// For example, you might extract a function similar to this:
//
//   func processBatch(ctx context.Context, batchID, totalBatches int, records []*database.Record, jsClient JetStreamPublisher) error { … }
//
// where JetStreamPublisher is an interface you define with a Publish method.
// This lets you swap in a fake implementation during testing.

//
// Here’s an example of a refactored batch processing function and its test.

// Define an interface so that we can pass a fake in tests.
type JetStreamPublisher interface {
	Publish(ctx context.Context, subject string, data []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

// processBatch contains the core logic from the goroutine in main.
func processBatch(ctx context.Context, batchID, totalBatches int, records []*database.Record, jsClient JetStreamPublisher) error {
	// In the original code, the DB query uses this formula:
	//   mod((id-1), totalBatches)+1 == batchID
	// For testing purposes, we simulate filtering the records.
	var filtered []*database.Record
	for _, r := range records {
		if ((r.ID-1)%totalBatches)+1 == batchID {
			filtered = append(filtered, r)
		}
	}

	data := dbRecordsToBytes(filtered)
	subject := fmt.Sprintf("records.%d", batchID)
	_, err := jsClient.Publish(ctx, subject, data)
	return err
}

// fakeJetStreamClient is a fake implementation of JetStreamPublisher for testing.
type fakeJetStreamClient struct {
	// published maps subject to published data.
	published map[string][]byte
}

func (f *fakeJetStreamClient) Publish(ctx context.Context, subject string, data []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	if f.published == nil {
		f.published = make(map[string][]byte)
	}
	f.published[subject] = data
	return &jetstream.PubAck{
		Stream:   "test-stream",
		Sequence: uint64(len(f.published)),
	}, nil
}

func TestProcessBatch(t *testing.T) {
	totalBatches := 3
	// Create records with IDs 1 through 6.
	var records []*database.Record
	now := time.Now().UTC()
	for i := 1; i <= 6; i++ {
		records = append(records, &database.Record{
			ID:         i,
			InsertedAt: now,
		})
	}

	jsClient := &fakeJetStreamClient{}

	// Test for batchID = 2.
	err := processBatch(context.Background(), 2, totalBatches, records, jsClient)
	if err != nil {
		t.Fatalf("processBatch returned error: %v", err)
	}

	// For batchID = 2, the filtering condition is:
	//   ((ID - 1) mod 3) + 1 == 2
	// This should select records with IDs 2 and 5.
	subject := "records.2"
	publishedData, ok := jsClient.published[subject]
	if !ok {
		t.Fatalf("No message published for subject %s", subject)
	}

	var out []types.Record
	err = json.Unmarshal(publishedData, &out)
	if err != nil {
		t.Fatalf("Failed to unmarshal published data: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("Expected 2 records in published message, got %d", len(out))
	}
	// Verify the IDs of the published records.
	expectedIDs := []int{2, 5}
	for i, rec := range out {
		if rec.ID != expectedIDs[i] {
			t.Errorf("Expected record ID %d, got %d", expectedIDs[i], rec.ID)
		}
	}
}
