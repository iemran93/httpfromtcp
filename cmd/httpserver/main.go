package main

import (
	"fmt"
	"learnhttp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const PORT = 42069

func main() {
	// TODO: implement handler function
	server, err := server.Serve(PORT)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", fmt.Sprintf(":%d", PORT))

	// server running until CTRL+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
