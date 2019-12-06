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
	Reset          bool
}

func loadEnv() (cfg config) {
	// Assign config keys
	cfg = config{
		Host:           os.Getenv("HOST"),
		Port:           os.Getenv("PORT"),
		Database:       os.Getenv("DATABASE"),
		FilesDirectory: os.Getenv("FILES_DIR"),
		Reset:          false,
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
		cfg.FilesDirectory = "./files"
	}
	if reset := os.Getenv("RESET"); reset == "YES" || reset == "yes" {
		cfg.Reset = true
	}

	// Set path as absolute
	cfg.FilesDirectory, _ = filepath.Abs(cfg.FilesDirectory)

	return
}
