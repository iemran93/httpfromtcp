package handler

import (
	"io"
	"learnhttp/internal/request"
)

type Handler func(w io.Writer, r *request.Request) *HandlerError

type HandlerError struct {
	StatusCode int
	Message    string
}
