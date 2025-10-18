package server

import (
	"fmt"
	"learnhttp/internal/handler"
	"learnhttp/internal/request"
	"learnhttp/internal/response"
	"net"
)

type Server struct {
	Listner net.Listener
	Handler handler.Handler
	State   ServerState
}

type ServerState int

const (
	ServerInit ServerState = iota
	serverClosed
)

func Serve(port int, handler handler.Handler) (*Server, error) {
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := Server{
		State:   ServerInit,
		Listner: listner,
		Handler: handler,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) listen() {
	// accept new connections
	for {
		c, err := s.Listner.Accept()

		if err != nil {
			if s.State == serverClosed {
				return
			}
			fmt.Printf("error accepting connection: %v", err)
			continue
		}

		go s.handle(c)
	}
}

func (s *Server) handle(c net.Conn) {
	// handle each connection
	defer c.Close()

	writer := response.NewWriter(c)
	request, err := request.RequestFromReader(c)
	if err != nil {
		he := &handler.HandlerError{
			StatusCode: 400,
			Message:    err.Error(),
		}
		he.WriteError(writer)
		return
	}

	// // handler write to here
	// var b bytes.Buffer
	// writer.Writer = &b

	// call the handler
	writer.Headers = make(map[string]string)
	s.Handler(writer, request)
}

func (s *Server) Close() error {
	// close the listner and server
	s.State = serverClosed
	return s.Listner.Close()
}
