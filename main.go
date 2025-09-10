package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChan := make(chan string)
	current := ""
	go func() {
		for {
			buff := make([]byte, 8)
			n, err := f.Read(buff)
			str := string(buff[:n])
			parts := strings.Split(str, "\n")
			for _, l := range parts[:len(parts)-1] {
				lineChan <- current + l
				current = ""
			}
			current += parts[len(parts)-1]
			if err != nil {
				if current != "" {
					lineChan <- current
					current = ""
				}
				if errors.Is(err, io.EOF) {
					close(lineChan)
					f.Close()
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				close(lineChan)
				f.Close()
				break
			}
		}
	}()

	return lineChan
}

func main() {
	tcpListener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer tcpListener.Close()

	for {
		connect, err := tcpListener.Accept()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Connection has been accepted")

		lines := getLinesChannel(connect)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("Connection has been closed")
	}
}
