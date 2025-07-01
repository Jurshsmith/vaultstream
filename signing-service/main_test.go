// main_test.go
package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/jurshsmith/vaultstream/database"
	"github.com/jurshsmith/vaultstream/types"
)

// ----------------------------
// Tests for Signing Logic
// ----------------------------

// TestSignRecord verifies that signRecord produces the expected signature.
func TestSignRecord(t *testing.T) {
	// Arrange: set up a record and a key.
	rec := types.Record{ID: 1}
	key := types.Key{ID: 10, Value: "secret-key"}

	// Act: sign the record.
	sig := signRecord(rec, &key)

	// Compute expected signature value.
	data := fmt.Sprintf("%d:%s", rec.ID, key.Value)
	hash := sha256.Sum256([]byte(data))
	expectedValue := base64.StdEncoding.EncodeToString(hash[:])

	// Assert: verify that all fields match.
	if sig.RecordID != rec.ID {
		t.Errorf("Expected RecordID %d, got %d", rec.ID, sig.RecordID)
	}
	if sig.KeyID != key.ID {
		t.Errorf("Expected KeyID %d, got %d", key.ID, sig.KeyID)
	}
	if sig.Value != expectedValue {
		t.Errorf("Expected signature value %q, got %q", expectedValue, sig.Value)
	}
}

// TestSignRecords verifies that signRecords concurrently signs all records correctly
// and preserves the input order.
func TestSignRecords(t *testing.T) {
	// Arrange: create several records.
	records := []types.Record{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	key := types.Key{ID: 5, Value: "test-key"}

	// Act: sign all records.
	sigs, err := signRecords(records, &key)
	if err != nil {
		t.Fatalf("signRecords returned an unexpected error: %v", err)
	}

	// Assert: the number of signatures must equal the number of records.
	if len(sigs) != len(records) {
		t.Errorf("Expected %d signatures, got %d", len(records), len(sigs))
	}

	// Check each signature is as expected and in the same order as the input records.
	for i, rec := range records {
		expected := signRecord(rec, &key)
		if sigs[i].RecordID != expected.RecordID ||
			sigs[i].KeyID != expected.KeyID ||
			sigs[i].Value != expected.Value {
			t.Errorf("Signature mismatch for record %d. Expected %+v, got %+v", rec.ID, expected, sigs[i])
		}
	}
}

// ----------------------------
// Integration Tests for Bulk Insertion Using Actual Database Connection (ent)
// ----------------------------

// setupDB ensures that both the signatures and records tables are clean,
// and inserts dummy records with IDs 1, 2, and 3 so that foreign key constraints pass.
func setupDB(t *testing.T) (*database.Client, context.Context) {
	dbClient := database.Connect()
	ctx := context.Background()

	// Clean the signatures table.
	if _, err := dbClient.Signature.Delete().Exec(ctx); err != nil {
		t.Skipf("Skipping integration test: unable to clean signatures table: %v", err)
	}

	// Clean the records table.
	if _, err := dbClient.Record.Delete().Exec(ctx); err != nil {
		t.Skipf("Skipping integration test: unable to clean records table: %v", err)
	}

	// Insert dummy records with IDs 1, 2, and 3.
	// Adjust the creation logic if your ent schema disallows manual ID setting.
	if _, err := dbClient.Record.Create().SetID(1).Save(ctx); err != nil {
		t.Skipf("Skipping integration test: unable to insert dummy record 1: %v", err)
	}
	if _, err := dbClient.Record.Create().SetID(2).Save(ctx); err != nil {
		t.Skipf("Skipping integration test: unable to insert dummy record 2: %v", err)
	}
	if _, err := dbClient.Record.Create().SetID(3).Save(ctx); err != nil {
		t.Skipf("Skipping integration test: unable to insert dummy record 3: %v", err)
	}

	return dbClient, ctx
}

// TestInsertSignaturesSuccess verifies that insertSignatures returns no error
// and that the signatures table contains the expected records after insertion.
func TestInsertSignaturesSuccess(t *testing.T) {
	// Connect to the actual database.
	dbClient, ctx := setupDB(t)
	defer dbClient.Close()

	// Create a small list of signatures (assumed to be fewer than the bulk chunk size).
	sigs := []types.Signature{
		{RecordID: 1, KeyID: 10, Value: "sig1"},
		{RecordID: 2, KeyID: 10, Value: "sig2"},
		{RecordID: 3, KeyID: 10, Value: "sig3"},
	}

	// Act: attempt to insert the signatures.
	if err := insertSignatures(ctx, dbClient, sigs); err != nil {
		t.Fatalf("insertSignatures returned error: %v", err)
	}

	// Give the DB a moment if needed.
	time.Sleep(100 * time.Millisecond)

	// Query the signatures table to verify the insert.
	actualSigs, err := dbClient.Signature.Query().All(ctx)
	if err != nil {
		t.Fatalf("failed to query signatures: %v", err)
	}

	if len(actualSigs) != len(sigs) {
		t.Errorf("Expected %d signatures in DB, got %d", len(sigs), len(actualSigs))
	}

	// Verify that each inserted signature matches the expected values.
	for _, expected := range sigs {
		found := false
		for _, actual := range actualSigs {
			// Adjust field access if your ent schema uses different names or types.
			if actual.RecordID == expected.RecordID &&
				actual.KeyID == expected.KeyID &&
				actual.Value == expected.Value {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected signature %+v not found in DB", expected)
		}
	}
}

// TestInsertSignaturesFailure simulates a failure during bulk insert by using a canceled context.
func TestInsertSignaturesFailure(t *testing.T) {
	// Connect to the actual database.
	dbClient, ctx := setupDB(t)
	defer dbClient.Close()

	sigs := []types.Signature{
		{RecordID: 1, KeyID: 10, Value: "sig1"},
	}

	// Create a context that is already canceled to simulate failure.
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // cancel immediately

	// Act: attempt to insert the signatures with the canceled context.
	err := insertSignatures(cancelCtx, dbClient, sigs)
	if err == nil {
		t.Fatal("Expected insertSignatures to return an error due to canceled context, but got nil")
	}

	// Verify that no signatures were inserted.
	actualSigs, err := dbClient.Signature.Query().All(ctx)
	if err != nil {
		t.Fatalf("failed to query signatures: %v", err)
	}
	if len(actualSigs) != 0 {
		t.Errorf("Expected 0 signatures in DB after failure, got %d", len(actualSigs))
	}
}
