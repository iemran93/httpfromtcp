package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(r io.Reader) (*Request, error) {
	request := Request{}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\r\n")
	rl, err := parseRequestLine(lines[0])
	if err != nil {
		return &request, err
	}

	request.RequestLine = *rl

	return &request, nil
}

func parseRequestLine(rl string) (*RequestLine, error) {
	requestLine := RequestLine{}
	rlp := strings.Split(rl, " ")
	if len(rlp) != 3 {
		return &requestLine, errors.New("Request line len error")
	}
	for i, p := range rlp {
		switch i {
		// method
		case 0:
			if p == "POST" || p == "GET" {
				requestLine.Method = p
			} else {
				return &requestLine, errors.New("Wrong http method!")
			}
		// path
		case 1:
			requestLine.RequestTarget = p
		// http version
		case 2:
			httpVersion := strings.Split(p, "/")
			if len(httpVersion) != 2 {
				return &requestLine, errors.New("Invalid http version!")
			}
			if httpVersion[1] == "1.1" {
				requestLine.HttpVersion = httpVersion[1]
			} else {
				return &requestLine, errors.New("Invalid http version!")
			}
		}
	}

	return &requestLine, nil
}
