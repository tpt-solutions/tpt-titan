package main

import (
	"log"
	"tpt-titan-simple/backend/config"
	"tpt-titan-simple/backend/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create and initialize server
	srv := server.NewServer(cfg)
	if err := srv.Initialize(); err != nil {
		log.Fatal("Failed to initialize server:", err)
	}
	defer srv.Close()

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
