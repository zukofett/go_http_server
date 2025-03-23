package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/textproto"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type config struct {
	port  int
	proto string
}

type Server struct {
	logger *slog.Logger
	config config
}

type Header map[string][]string

type RequestLine struct {
	Method     string
	RequestURL *url.URL
	Proto      string
	ProtoMajor int
	ProtoMinor int
}

type Request struct {
	RequestLine
	Header
	RequestURI string
	Body       io.ReadCloser
}

func NewServer(port int, proto string) *Server {
	return &Server{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		config: config{
			port:  port,
			proto: proto,
		},
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen(s.config.proto, fmt.Sprintf(":%d", s.config.port))
	if err != nil {
		return err
	}

	return s.Serve(listener)
}

func (s *Server) Serve(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}

		go s.handleRequest(conn)
	}
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	req, err := readRequest(r)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}

	fmt.Printf("%+v\n", req)
}

func readRequest(r *bufio.Reader) (*Request, error) {
	tpReader := textproto.NewReader(r)

	req := new(Request)

	var str string
	var err error

	str, err = tpReader.ReadLine()
	if err != nil {
		return nil, err
	}

	var ok bool
	req.Method, req.RequestURI, req.Proto, ok = parseRequestLine(str)
	if !ok {
		return nil, fmt.Errorf("Bad request line %q", str)
	}

	if !isValidMethod(req.Method) {
		return nil, fmt.Errorf("invalid method: %q", req.Method)
	}

	req.RequestURL, err = url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return nil, err
	}

    req.ProtoMajor, req.ProtoMinor, ok = parseProtoVersion(req.Proto)
	if !ok {
        return nil, fmt.Errorf("Invalid proto: %q", req.Proto)
	}

    mimeHeaders, err := tpReader.ReadMIMEHeader()
    if err != nil {
        return nil, err
    }

    req.Header = Header(mimeHeaders)

    // TODO: read body

    return req, nil
}

func parseRequestLine(line string) (string, string, string, bool) {
	method, rest, found1 := strings.Cut(line, " ")
	url, proto, found2 := strings.Cut(rest, " ")
	if !found1 || !found2 {
		return "", "", "", false
	}

	return method, url, proto, true
}

func parseProtoVersion(proto string) (int, int, bool) {
	if len(proto) != len("HTTP/X.X") {
		return 0, 0, false
	}

	if !strings.HasPrefix(proto, "HTTP/") {
		return 0, 0, false
	}

	version := strings.TrimPrefix(proto, "HTTP/")
	nums := strings.Split(version, ".")
	if len(nums) != 2 {
		return 0, 0, false
	}

	major, err := strconv.Atoi(nums[0])
	if err != nil {
		return 0, 0, false
	}

	minor, err := strconv.Atoi(nums[1])
	if err != nil {
		return 0, 0, false
	}

    return major, minor, true
}

func main() {
	server := NewServer(4000, "tcp")

	err := server.ListenAndServe()
	server.logger.Error(err.Error())
	os.Exit(1)
}
