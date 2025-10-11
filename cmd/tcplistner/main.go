package main

import (
	"fmt"
	"learnhttp/internal/request"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			req, err := request.RequestFromReader(c)
			if err != nil {
				log.Fatalf("exit with error: %v", err)
			}

			fmt.Printf(`
Request line:
- Method: %v
- Target: %v
- Version: %v
Headers:
`, req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

			for k, v := range req.Headers {
				fmt.Printf("- %v: %v\n", k, v)
			}

			c.Close()
		}(c)

	}

}
