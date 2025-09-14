package main

import (
	"fmt"
	"log"
	"net"

	"github.com/zukofett/http/internal/request"
)

func main() {
	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("failed to open file", err)
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal("error in connection accept", err)
		}

		req, err := request.FromReader(conn)
		if err != nil {
			log.Fatal("error reading from the connection", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
	}

}
