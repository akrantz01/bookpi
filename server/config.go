package main

import (
	"os"
	"path/filepath"
)

type config struct {
	Host           string
	Port           string
	Database       string
	FilesDirectory string
}

func loadEnv() (cfg config) {
	// Assign config keys
	cfg = config{
		Host:           os.Getenv("HOST"),
		Port:           os.Getenv("PORT"),
		Database:       os.Getenv("DATABASE"),
		FilesDirectory: os.Getenv("FILES_DIR"),
	}

	// Set defaults if not exist
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.Database == "" {
		cfg.Database = "./database.db"
	}
	if cfg.FilesDirectory == "" {
		cfg.Database = "./files"
	}

	// Set path as absolute
	cfg.FilesDirectory, _ = filepath.Abs(cfg.FilesDirectory)

	return
}
