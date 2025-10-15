package server

import (
	"bytes"
	"fmt"
	"learnhttp/internal/handler"
	"learnhttp/internal/request"
	"learnhttp/internal/response"
	"log/slog"
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

		go func() {
			slog.Info("#go-con", "new con", c)
			s.handle(c)
		}()
	}
}

func (s *Server) handle(c net.Conn) {
	// handle each connection
	defer c.Close()

	request, err := request.RequestFromReader(c)
	slog.Info("#handle", "request", request)
	if err != nil {
		response.WriteError(
			c,
			&handler.HandlerError{
				StatusCode: 400,
				Message:    err.Error(),
			})
		return
	}

	// handler write to here
	var b bytes.Buffer

	// call the handler
	handlerError := s.Handler(&b, request)
	if handlerError != nil { // send handler error
		err := response.WriteError(c, handlerError)
		if err != nil {
			return
		}
		return
	}

	h := response.GetDefaultHeaders(b.Len())
	err = response.WriteStatusLine(c, response.Ok)
	if err != nil {
		return
	}
	response.WriteHeaders(c, h)
	c.Write(b.Bytes())
	// c.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"))
	// c.Close()
	// err := response.WriteStatusLine(c, response.Ok)
	// if err != nil {
	// 	return
	// }
	// h := response.GetDefaultHeaders(0)
	// response.WrtieHeaders(c, h)
	// c.Close()
}

func (s *Server) Close() error {
	// close the listner and server
	s.State = serverClosed
	return s.Listner.Close()
}
