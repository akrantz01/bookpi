package main

import "os"

type config struct {
	Host     string
	Port     string
	Database string
}

func loadEnv() (cfg config) {
	// Assign config keys
	cfg = config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Database: os.Getenv("DATABASE"),
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

	return
}
