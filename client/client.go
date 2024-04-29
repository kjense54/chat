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
	
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		input = strings.TrimSpace(input)

		// quit client
		if input == "/quit" {
			fmt.Println("Quitting")
			return
		} else {
			// write test message
			if _, err := conn.Write([]byte(string(len(input)) + input)); err != nil {
				log.Fatal(err)
			}
		}
	}
}
