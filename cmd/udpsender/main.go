package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("error in resolving address", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("error openning connection:", err)
	}
    defer func() {
        if err := conn.Close(); err != nil {
            log.Fatal("err closing file:", err)
        }
    }()

	reader := bufio.NewReader(os.Stdin)
	for {
        fmt.Print("> ")
		
        line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error reading line:", err)
		}
        
        _, err = conn.Write([]byte(line))
        if err != nil {
			log.Fatal("error writing line:", err)
		}
	}
}
