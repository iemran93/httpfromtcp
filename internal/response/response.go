package response

import (
	"fmt"
	"io"
	"learnhttp/internal/handler"
	"learnhttp/internal/headers"
	"strconv"
)

type StatusCode int

const (
	Ok          StatusCode = 200
	ClientError StatusCode = 400
	ServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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
	_, err := w.Write(fmt.Appendf(nil, "%v\r\n", statuseLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for k, v := range h {
		_, err := w.Write(fmt.Appendf(nil, "%v: %v\r\n", k, v))
		if err != nil {
			return err
		}

	}
	_, err := w.Write([]byte("\r\n"))
	return err
}

func WriteError(w io.Writer, handlerError *handler.HandlerError) error {
	contentLen := len(handlerError.Message)
	err := WriteStatusLine(w, StatusCode(handlerError.StatusCode))
	if err != nil {
		return err
	}

	h := GetDefaultHeaders(contentLen)
	WriteHeaders(w, h)
	w.Write([]byte(handlerError.Message))

	return nil
}
