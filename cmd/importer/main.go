package main

import (
	"context"
	"log"
	"os"

	"tool-service/internal/config"
	"tool-service/internal/importer"
	"tool-service/internal/repository"
	phonerepo "tool-service/internal/repository/phone"
)

func main() {
	cfg, err := config.LoadImporter()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	ctx := context.Background()

	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	if err := repository.Migrate(ctx, db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	rows, err := importer.ImportMobileCSV(cfg.ImporterCSVPath)
	if err != nil {
		log.Fatalf("csv: %v", err)
	}

	repo := phonerepo.New(db)
	if err := repo.Truncate(ctx); err != nil {
		log.Fatalf("truncate: %v", err)
	}
	if err := repo.InsertBatch(ctx, rows); err != nil {
		log.Fatalf("insert: %v", err)
	}
	log.Printf("imported %d rows from %s", len(rows), cfg.ImporterCSVPath)
	os.Exit(0)
}
