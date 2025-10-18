package request

import (
	"errors"
	"fmt"
	"io"
	"learnhttp/internal/headers"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
	ParsingHeaders
	ParsingBody
)

func (r *Request) parse(data []byte) (int, error) {
	totalByteConsumed := 0
Outer:
	for {
		switch r.ParserState {
		case ParserInit:
			// parse request line
			rl, bc, err := parseRequestLine(data[totalByteConsumed:])
			if err != nil {
				return 0, err
			} else if bc == 0 {
				// zero bytes consumed; need more data
				break Outer
			}
			r.RequestLine = *rl
			r.ParserState = ParsingHeaders
			totalByteConsumed += bc
		case ParsingHeaders:
			// parse headers
			n, done, err := r.Headers.Parse(data[totalByteConsumed:])
			if err != nil {
				return 0, err
			}
			totalByteConsumed += n
			if done {
				if _, exist := r.Headers.Get("content-length"); !exist {
					r.ParserState = ParserDone
				} else {
					r.ParserState = ParsingBody
				}
				// the end of headers \r\n
				totalByteConsumed += 2
				break Outer
			}

			if n == 0 {
				// need more data
				break Outer
			}
		case ParsingBody:
			// parse body
			contLengthS, exist := r.Headers.Get("content-length")
			if !exist {
				r.ParserState = ParserDone
				break Outer
			}
			contLength, err := strconv.Atoi(contLengthS)
			if err != nil {
				return 0, err
			}

			if contLength == 0 {
				r.ParserState = ParserDone
				break Outer
			}

			remaining := min(contLength-len(r.Body), len(data[totalByteConsumed:]))
			r.Body = append(r.Body, data[totalByteConsumed:totalByteConsumed+remaining]...)
			if len(r.Body) == contLength {
				r.ParserState = ParserDone
			}
			totalByteConsumed += remaining
			break Outer
		case ParserDone:
			return 0, errors.New("error: trying to read data in a done state")
		default:
			return 0, errors.New("error: unknown state")
		}
	}
	return totalByteConsumed, nil
}

func RequestFromReader(r io.Reader) (*Request, error) {
	request := Request{
		ParserState: ParserInit,
		RequestLine: RequestLine{},
		Headers:     headers.NewHeaders(),
		Body:        make([]byte, 0),
	}

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

	// expect more data for the body
	if request.ParserState == ParsingBody {
		contLengthS, exist := request.Headers.Get("content-length")
		if exist {
			contLength, _ := strconv.Atoi(contLengthS)
			if len(request.Body) < contLength {
				return nil, fmt.Errorf("incomplete body: expected %v bytes, got %v", contLength, len(request.Body))
			}
		}
	}
	request.ParserState = ParserDone
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

	return &requestLine, rlb + 2, nil
}
