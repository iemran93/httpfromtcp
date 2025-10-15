package server

import (
	"fmt"
	"learnhttp/internal/response"
	"net"
)

type Server struct {
	Listner net.Listener
	State   ServerState
}

type ServerState int

const (
	ServerInit ServerState = iota
	serverClosed
)

func Serve(port int) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := Server{State: ServerInit, Listner: listner}
	// defer listner.Close()

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	// close the listner and server
	s.State = serverClosed
	return s.Listner.Close()
}

func (s *Server) listen() {
	// accept new connections
	for {
		c, err := s.Listner.Accept()
		if s.State == serverClosed {
			return
		}

		if err != nil {
			return
		}

		go s.handle(c)
	}
}

func (s *Server) handle(c net.Conn) {
	// handle each connection
	// _, err := request.RequestFromReader(c)
	// if err != nil {
	// 	return
	// }
	// c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"))
	// c.Close()
	err := response.WriteStatusLine(c, response.Ok)
	if err != nil {
		return
	}
	h := response.GetDefaultHeaders(0)
	response.WrtieHeaders(c, h)
	c.Close()
}
