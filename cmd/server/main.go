package main

import (
	"log"
)

func main() {
	// Create and start server
	server := NewServer()

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
