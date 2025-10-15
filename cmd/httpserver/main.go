package main

import (
	"fmt"
	"io"
	"learnhttp/internal/handler"
	"learnhttp/internal/request"
	"learnhttp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const PORT = 42069

func testHandler(w io.Writer, r *request.Request) *handler.HandlerError {
	path := r.RequestLine.RequestTarget
	if path == "/yourproblem" {
		return &handler.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}
	} else if path == "/myproblem" {
		return &handler.HandlerError{
			StatusCode: 500,
			Message:    "Oops, my bad\n",
		}
	}
	w.Write([]byte("All good\n"))
	return nil
}
func main() {
	server, err := server.Serve(PORT, testHandler)
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
