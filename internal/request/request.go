package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

var methodRX = regexp.MustCompile("^[A-Z]+$")

var ErrorBadRequestLine = fmt.Errorf("malformed request line")
var ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http version")
var ErrorRequestInErrorState = fmt.Errorf("request is in error state")

var SEPARATOR = []byte{'\r', '\n'}

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func FromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	readIdx := 0
	for !request.done() {
		n, err := reader.Read(buf[readIdx:])
		if err != nil {
			return nil, err
		}

		readIdx += n

		read, err := request.parse(buf[:readIdx+n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[read:readIdx])
		readIdx -= read
	}

	return request, nil
}

func (r *Request) parse(msg []byte) (int, error) {
	read := 0

	for {
		switch r.state {
		case StateInit:
            rl, n, err := parseRequestLine(msg[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				return read, nil
			}

			r.RequestLine = *rl
			read += n

			r.state = StateDone

		case StateDone:
			return read, nil
		case StateError:
			return 0, ErrorRequestInErrorState
		}
	}
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func parseRequestLine(msg []byte) (*RequestLine, int, error) {
	idx := bytes.Index(msg, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := msg[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Fields(startLine)
	if len(parts) != 3 {
		return nil, 0, ErrorBadRequestLine
	}

	if !methodRX.Match(parts[0]) {
		return nil, 0, ErrorBadRequestLine
	}

	version := bytes.Split(parts[2], []byte("/"))
	if len(version) != 2 || string(version[0]) != "HTTP" || string(version[1]) != "1.1" {
		return nil, 0, ErrorUnsupportedHttpVersion
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(version[1]),
	}

	return rl, read, nil
}
