package response

import (
	"fmt"
	"io"
	"learnhttp/internal/headers"
	"maps"
)

type Writer struct {
	Writer  io.Writer
	Headers map[string]string
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: w,
	}
}

type StatusCode int

const (
	Ok          StatusCode = 200
	ClientError StatusCode = 400
	ServerError StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	// map status code to a reason phrase
	statuseLine := fmt.Sprintf("HTTP/1.1 %v ", statusCode)
	switch statusCode {
	case Ok:
		statuseLine += "OK"
	case ClientError:
		statuseLine += "Bad Request"
	case ServerError:
		statuseLine += "Intenal Server Error"
	default:
	}
	_, err := w.Writer.Write(fmt.Appendf(nil, "%v\r\n", statuseLine))
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	for k, v := range h {
		_, err := w.Writer.Write(fmt.Appendf(nil, "%v: %v\r\n", k, v))
		if err != nil {
			return err
		}

	}
	_, err := w.Writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (w *Writer) WriteChunkedBody(b []byte) (int, error) {
	// write the chunked format
	// hex representation of an int %x base 16
	out := fmt.Sprintf("%x\r\n", len(b))
	n, err := w.WriteBody([]byte(out))
	if err != nil {
		return 0, err
	}

	n, err = w.WriteBody(b)
	if err != nil {
		return 0, err
	}

	n, err = w.WriteBody([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteChunckedBodyDone() (int, error) {
	out := fmt.Sprintf("%x\r\n", 0)
	n, err := w.WriteBody([]byte(out))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func GetHeaders(newHeaders map[string]string) headers.Headers {
	h := headers.NewHeaders()
	// deafult
	h["Content-Length"] = "0"
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	maps.Copy(h, newHeaders)
	return h
}
