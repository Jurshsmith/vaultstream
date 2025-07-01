package main

import (
	"context"
	"fmt"

	"github.com/jurshsmith/vaultstream/config"
	db "github.com/jurshsmith/vaultstream/database"
	vaultStreamLogger "github.com/jurshsmith/vaultstream/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	logger = vaultStreamLogger.New()
	defer logger.Sync()

	logger.Info("Starting the seeding process.")

	config.Setup()

	totalRecords := config.TotalRecords()
	logger.Info("Total records to seed", zap.Int("totalRecords", totalRecords))

	client := db.Connect()
	defer client.Close()

	ctx := context.Background()

	if err := clearTable(ctx, client, "signatures"); err != nil {
		logger.Debug("Failed to clear signatures table", zap.Error(err))
	}

	if err := clearTable(ctx, client, "records"); err != nil {
		logger.Debug("Failed to clear records table", zap.Error(err))
	}

	if err := seedRecords(ctx, client, totalRecords); err != nil {
		logger.Fatal("Failed seeding records", zap.Error(err))
	}

	logger.Info("Seeding records complete!")
}

func clearTable(ctx context.Context, client *db.Client, table string) error {
	logger.Info("Clearing table", zap.String("table", table))
	_, err := client.Exec(ctx, fmt.Sprintf("DELETE FROM %s;", table))
	if err != nil {
		logger.Error("Error clearing table", zap.String("table", table), zap.Error(err))
	}
	return err
}

func seedRecords(ctx context.Context, client *db.Client, totalRecords int) error {
	logger.Info("Seeding records", zap.Int("totalRecords", totalRecords))
	query := fmt.Sprintf(`
		INSERT INTO records (inserted_at)
		SELECT now()
		FROM generate_series(1, %d);
	`, totalRecords)

	_, err := client.Exec(ctx, query)
	if err != nil {
		logger.Error("Error seeding records", zap.Int("totalRecords", totalRecords), zap.Error(err))
	}
	return err
}
