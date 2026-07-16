package main

import (
	"log"
	"os"
	"tpt-titan/backend/cmd"
	"tpt-titan/backend/config"
	"tpt-titan/backend/internal/server"
)

func main() {
	// If the first argument is a known management subcommand, run it and exit.
	if cmd.Run(os.Args[1:]) {
		return
	}

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
