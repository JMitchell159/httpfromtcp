package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/JMitchell159/httpfromtcp/internal/request"
)

func main() {
	tcpListener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer tcpListener.Close()

	for {
		connect, err := tcpListener.Accept()
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		fmt.Println("Connection has been accepted")

		req, err := request.RequestFromReader(connect)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("Connection has been closed")
	}
}
