package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	ParserState ParserState
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type ParserState int

const (
	ParserInit ParserState = iota
	ParserDone
)

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParserState {
	case ParserInit:
		rl, bc, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		} else if bc == 0 {
			// zero bytes consumed; need more data
			return 0, nil
		}
		r.RequestLine = *rl
		r.ParserState = ParserDone
		return bc, nil
	case ParserDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}

func RequestFromReader(r io.Reader) (*Request, error) {
	request := Request{ParserState: ParserInit, RequestLine: RequestLine{}}

	bufferSize := 8
	buf := make([]byte, bufferSize)
	readToIndx := 0
	for request.ParserState != ParserDone {
		// instead of appending bytes to buffer
		// if bytes consumed shift to that index

		// if buffer full extend it
		if readToIndx == cap(buf) {
			bufferSize *= 2
			newBuf := make([]byte, bufferSize)
			copy(newBuf, buf)
			buf = newBuf
		}

		byte_read, err := r.Read(buf[readToIndx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.ParserState = ParserDone
				break
			}
			return nil, err
		}

		readToIndx += byte_read

		byte_parse, err := request.parse(buf[:readToIndx])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[byte_parse:readToIndx])
		readToIndx -= byte_parse
	}

	return &request, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	requestLine := RequestLine{}

	lines_s := strings.Split(string(b), "\r\n")

	if len(lines_s) == 1 {
		// still not reach \r\n; return 0 no error
		return nil, 0, nil
	}

	rlb := len(lines_s[0])
	rlp := strings.Split(lines_s[0], " ")
	if len(rlp) != 3 {
		return &requestLine, rlb, errors.New("Request line len error")
	}
	for i, p := range rlp {
		switch i {
		// method
		case 0:
			if p == "POST" || p == "GET" {
				requestLine.Method = p
			} else {
				return &requestLine, rlb, errors.New("Wrong http method!")
			}
		// path
		case 1:
			requestLine.RequestTarget = p
		// http version
		case 2:
			httpVersion := strings.Split(p, "/")
			if len(httpVersion) != 2 {
				return &requestLine, rlb, errors.New("Invalid http version!")
			}
			if httpVersion[1] == "1.1" {
				requestLine.HttpVersion = httpVersion[1]
			} else {
				return &requestLine, rlb, errors.New("Invalid http version!")
			}
		}
	}

	return &requestLine, rlb, nil
}
