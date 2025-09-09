package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Printf("error when opening messages.txt: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	current := ""
	for {
		buff := make([]byte, 8)
		n, err := file.Read(buff)
		str := string(buff[:n])
		parts := strings.Split(str, "\n")
		for _, l := range parts[:len(parts)-1] {
			fmt.Printf("read: %s%s\n", current, l)
			current = ""
		}
		current += parts[len(parts)-1]
		if err != nil {
			if current != "" {
				fmt.Printf("read: %s\n", current)
				current = ""
			}
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
	}
}
