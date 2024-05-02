package main

import (
	"context"
	"log"
	"net"
	"time"	
	"fmt"
	"bufio"
	"os"
	"strings"
	"io"
)

func main() {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", "localhost:8090")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	
	inputCh := make(chan string)
	// Write to server
	go func() {
		var input string
		for {
			reader := bufio.NewReader(os.Stdin)
			input, err = reader.ReadString('\n')
			if err != nil {
			log.Fatalf("Error reading input: %v", err)
			}
			input = strings.TrimSpace(input)
			inputCh <- input
		}
	}()

	// Read from server
	go func() {
		defer conn.Close()
		
		buff := make([]byte, 106) // largest msg size
		c := bufio.NewReader(conn)
		
		for {
			// read first byte for message length
			size, err := c.ReadByte()
			if err != nil {
				return
			}
			// read full message to buff
			_, err1 := io.ReadFull(c, buff[:int(size)])
			if err1 != nil {
				return
			}
			// print
			fmt.Println(string(buff))
		}
	}()

	for {
		select {
		case input := <-inputCh:
			if input == "/quit" {
				fmt.Println("Quitting Client")
				conn.Close()
				return
			} else {
				// send message to server
				if _, err := conn.Write([]byte(string(len(input)) + input)); err != nil {
					log.Fatalf("Error writing to server: %v", err)
				}
			}
		}
		
	}
}
