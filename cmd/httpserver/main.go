package main

import (
	"fmt"
	"io"
	"learnhttp/internal/request"
	"learnhttp/internal/response"
	"learnhttp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const PORT = 42069

func testHandler(w *response.Writer, r *request.Request) {
	path := r.RequestLine.RequestTarget
	var statusCode response.StatusCode
	var body string

	w.Headers["Content-Type"] = "text/plain"
	h := response.GetHeaders(w.Headers)

	switch {
	case path == "/yourproblem":
		statusCode = 400
		body = response400()
	case path == "/myproblem":
		statusCode = 500
		body = response500()
	case strings.HasPrefix(path, "/httpbin/"):
		h["Transfer-Encoding"] = "chunked"
		delete(h, "Content-Length")
		target := path[len("/httpbin/"):]

		resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", target))
		if err != nil {
			statusCode = 500
			body = response500()
		}

		w.WriteStatusLine(response.Ok)
		w.WriteHeaders(h)
		for {
			buf := make([]byte, 32)
			n, err := resp.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					// done
					w.WriteChunckedBodyDone()
					break
				}
				statusCode = 500
				body = response500()
				break
			}
			log.Printf("Read from bin: \n%v\n%v\n", n, string(buf))
			w.WriteChunkedBody(buf[:n])
		}
		return
	default:
		statusCode = 200
		body = response200()
	}

	h["Content-Length"] = fmt.Sprintf("%d", len(body))

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(h)
	w.WriteBody([]byte(body))
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

func response200() string {
	return `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
}

func response400() string {
	return `
<html>
	<head>
	    <title>400 Bad Request</title>
	</head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
}

func response500() string {
	return `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
}
