package handler

import (
	"learnhttp/internal/request"
	"learnhttp/internal/response"
	"strconv"
)

type Handler func(w *response.Writer, r *request.Request) *HandlerError

type HandlerError struct {
	StatusCode int
	Message    string
}

func (he *HandlerError) WriteError(w *response.Writer) error {
	contentLen := len(he.Message)
	conentLenS := strconv.Itoa(contentLen)
	err := w.WriteStatusLine(response.StatusCode(he.StatusCode))
	if err != nil {
		return err
	}

	w.Headers["Content-Length"] = conentLenS
	h := response.GetHeaders(w.Headers)
	w.WriteHeaders(h)
	w.WriteBody([]byte(he.Message))

	return nil

}
