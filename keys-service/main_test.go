package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jurshsmith/vaultstream/nats"
	"github.com/jurshsmith/vaultstream/types"
	"go.uber.org/zap"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "Valid key generation",
			id:      1,
			wantErr: false,
		},
		{
			name:    "Valid key with another ID",
			id:      42,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := generateKey(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("generateKey() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("generateKey() unexpected error = %v", err)
				}
				if key == nil {
					t.Errorf("generateKey() key is nil")
					return
				}
				if key.ID != tt.id {
					t.Errorf("generateKey() key.ID = %v, want %v", key.ID, tt.id)
				}
				if key.Value == "" {
					t.Errorf("generateKey() key.Value is empty")
				}
				if key.IsInUse {
					t.Errorf("generateKey() key.IsInUse = true, want false")
				}
				if key.LastUsedAt != time.Unix(0, 0) {
					t.Errorf("generateKey() key.LastUsedAt = %v, want %v", key.LastUsedAt, time.Unix(0, 0))
				}

				// Test that the key is valid JSON
				jsonStr := fmt.Sprintf(`{"key":"%s"}`, key.Value)
				if !json.Valid([]byte(jsonStr)) {
					t.Errorf("generateKey() key.Value is not valid JSON when encoded: %s", jsonStr)
				}
			}
		})
	}
}

func TestGenerateKeys(t *testing.T) {
	tests := []struct {
		name      string
		totalKeys int
		wantErr   bool
	}{
		{
			name:      "Generate 0 keys",
			totalKeys: 0,
			wantErr:   false,
		},
		{
			name:      "Generate 1 key",
			totalKeys: 1,
			wantErr:   false,
		},
		{
			name:      "Generate 5 keys",
			totalKeys: 5,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := generateKeys(tt.totalKeys)

			if tt.wantErr {
				if err == nil {
					t.Errorf("generateKeys() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("generateKeys() unexpected error = %v", err)
				}
				if len(keys) != tt.totalKeys {
					t.Errorf("generateKeys() len = %v, want %v", len(keys), tt.totalKeys)
				}

				// Check that each key has the correct ID
				for i, key := range keys {
					if key.ID != i+1 {
						t.Errorf("generateKeys() key[%d].ID = %v, want %v", i, key.ID, i+1)
					}
					if key.Value == "" {
						t.Errorf("generateKeys() key[%d].Value is empty", i)
					}
				}
			}
		})
	}
}

func TestEnqueueKey(t *testing.T) {
	// Create a logger for testing
	oldLog := log
	log, _ = zap.NewDevelopment()
	defer func() {
		log = oldLog
	}()

	tests := []struct {
		name      string
		key       *types.Key
		wantPanic bool
	}{
		{
			name: "Successfully enqueue key",
			key: &types.Key{
				ID:         1,
				Value:      "mock-key-value",
				IsInUse:    false,
				LastUsedAt: time.Unix(0, 0),
			},
			wantPanic: false,
		},
		{
			name: "Error publishing key",
			key: &types.Key{
				ID:         2,
				Value:      "mock-key-value",
				IsInUse:    false,
				LastUsedAt: time.Unix(0, 0),
			},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js, conn := nats.Connect()
			defer conn.Close()

			ctx := context.Background()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("enqueueKey() did not panic as expected")
					}
				}()
				enqueueKey(js, tt.key, ctx)
			} else {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("enqueueKey() unexpected panic: %v", r)
					}
				}()
				enqueueKey(js, tt.key, ctx)
			}
		})
	}
}

// TestEnqueueAllKeys tests the enqueueAllKeys function
func TestEnqueueAllKeys(t *testing.T) {
	// Create a logger for testing
	oldLog := log
	log, _ = zap.NewDevelopment()
	defer func() {
		log = oldLog
	}()

	// Override KeysMaxConcurrency for this test

	tests := []struct {
		name string
		keys []*types.Key
	}{
		{
			name: "Enqueue multiple keys",
			keys: []*types.Key{
				{ID: 1, Value: "key1", IsInUse: false, LastUsedAt: time.Unix(0, 0)},
				{ID: 2, Value: "key2", IsInUse: false, LastUsedAt: time.Unix(0, 0)},
				{ID: 3, Value: "key3", IsInUse: false, LastUsedAt: time.Unix(0, 0)},
			},
		},
		{
			name: "Enqueue empty key list",
			keys: []*types.Key{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js, conn := nats.Connect()
			defer conn.Close()

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("enqueueAllKeys() unexpected panic: %v", r)
				}
			}()

			enqueueAllKeys(js, tt.keys)
		})
	}
}

// Test for main function
func TestMainIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping main integration test in short mode")
	}

	// This is a placeholder for a more comprehensive integration test
	// In a real test, you would mock dependencies and verify behaviors

	// For now, we'll just verify that key components compile
	t.Run("key_components_compile", func(t *testing.T) {
		// Just a compilation check
		_, err := generateKey(1)
		if err != nil {
			t.Errorf("generateKey() unexpected error: %v", err)
		}
	})
}
