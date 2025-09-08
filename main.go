package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error when reading messages.txt: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	buff := make([]byte, 8)
	for {
		if _, err = file.Read(buff); errors.Is(err, io.EOF) {
			os.Exit(0)
		}
		fmt.Printf("read: %s\n", buff)
		clear(buff)
	}
}
