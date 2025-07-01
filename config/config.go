package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Setup() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("Info: .env file not found, skipping preload")
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	} else {
		log.Println("Info: .env file successfully loaded")
	}
}

func EventsStreamName() string {
	return "vaultstream-streams"
}

func DatabaseURL() string {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("DATABASE_URL not set")
	}
	return url
}

func VaultStreamNatsURL() string {
	url := os.Getenv("VAULTSTREAM_NATS_URL")
	if url == "" {
		log.Fatal("VAULTSTREAM_NATS_URL not set")
	}
	return url
}
func VaultStreamNatsPassword() string {
	pwd := os.Getenv("VAULTSTREAM_NATS_PASSWORD")
	if pwd == "" {
		log.Fatal("VAULTSTREAM_NATS_PASSWORD not set")
	}
	return pwd
}

func TotalRecords() int {
	return mustEnvInt("TOTAL_RECORDS")
}
func TotalRecordBatches() int {
	return ceilDiv(TotalRecords(), RecordsBatchSize())
}

func KeysBucketName() string {
	return "vaultstream-keys"
}
func KeysTTLInSeconds() time.Duration {
	return time.Duration(100 * int(time.Second))
}
func TotalKeys() int {
	return mustEnvInt("TOTAL_KEYS")
}

func RecordsMaxConcurrency() int {
	return mustEnvInt("RECORDS_MAX_CONCURRENCY")
}
func RecordsBatchSize() int {
	return mustEnvInt("BATCH_SIZE")
}

func KeysMaxConcurrency() int {
	return mustEnvInt("KEYS_MAX_CONCURRENCY")
}

func SignerMaxConcurrency() int {
	return mustEnvInt("SIGNER_MAX_CONCURRENCY")
}

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

func mustEnvInt(key string) int {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("%s not set", key)
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return v
}
