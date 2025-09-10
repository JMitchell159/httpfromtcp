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
		log.Print(err)
		os.Exit(1)
	}

	connect, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer connect.Close()

	sc := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(" > ")
		line, err := sc.ReadString(byte('\n'))
		if err != nil {
			log.Print(err)
		}
		_, err = connect.Write([]byte(line))
		if err != nil {
			log.Print(err)
		}
	}
}
