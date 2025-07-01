package database

import (
	"log"

	vaultStreamConfig "github.com/jurshsmith/vaultstream/config"
	_ "github.com/lib/pq"
)

func Connect() *Client {
	dbURL := vaultStreamConfig.DatabaseURL()
	dbClient, err := Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed opening DB: %v", err)
	}
	return dbClient
}
