package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// convert string ip to UDP add type
	add, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		log.Fatal(err)
	}

	// laddress: listen
	// raddress: send
	c, err := net.DialUDP("udp", nil, add)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		_, err = c.Write([]byte(line))
		if err != nil {
			log.Println(err)
		}

	}

}

// ✗ go run ./cmd/udpsender/main.go
// > hi
// > hi
// 2025/10/01 21:29:31 write udp 127.0.0.1:34221->127.0.0.1:42069: write: connection refused

// ^ this error not related to udp connection
// the packet send the first time regarding if a conn is established
// this error comes from OS because it learn from first packet that conn not exist
