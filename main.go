package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// the io READET and WRITER interface work with any type of data stream(files, network)
func getLinesChannel(c io.ReadCloser) <-chan string {
	strs := make(chan string)
	go func() {
		defer close(strs)
		defer func() {
			fmt.Printf("connection closed %v\n", c)
			c.Close()
		}()

		c_line := ""
		for {
			b := make([]byte, 8)
			n, err := c.Read(b)
			if err != nil {
				break
			}

			parts := strings.Split(string(b[:n]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				c_line += parts[i]
				strs <- c_line
				c_line = ""
			}

			c_line += parts[len(parts)-1]
		}

		if c_line != "" {
			strs <- c_line
		}
		// close(strs)
	}()
	return strs
}

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

		for line := range getLinesChannel(c) {
			fmt.Printf("read: %v\n", line)
		}
	}

}
