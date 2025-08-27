package main

import (
	"fmt"
	"os"

	"github.com/adhitht/R2Mirror/internal/config"
	"github.com/adhitht/R2Mirror/internal/logger"
	"github.com/adhitht/R2Mirror/internal/processor"
	"github.com/adhitht/R2Mirror/internal/storage"
)

func main() {
	log := logger.New()
	log.Info("Ubuntu Release Downloader Starting...")

	if err := config.LoadEnv(); err != nil {
		log.Error("Environment error", "error", err)
		fmt.Println("💡 Please check your .env file and try again.")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error("Configuration error", "error", err)
		fmt.Println("💡 Please check your config.yaml file and try again.")
		os.Exit(1)
	}

	storageClient, err := storage.NewR2Client(cfg)
	if err != nil {
		log.Error("Storage client error", "error", err)
		fmt.Println("💡 Make sure your R2 credentials are properly configured.")
		os.Exit(1)
	}
	defer storageClient.Close()

	log.Info("Connected to R2", "bucket", cfg.Bucket, "region", cfg.Region)

	proc := processor.New(storageClient, log)
	
	log.Info("Starting initial processing...")
	if err := proc.ProcessReleases(cfg); err != nil {
		log.Error("Initial process failed", "error", err)
		fmt.Println("💡 Check your configuration and network connection.")
	}

	if err := proc.WatchConfig(cfg); err != nil {
		log.Error("Failed to start config watcher", "error", err)
		os.Exit(1)
	}
}