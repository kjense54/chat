package main

import (
	"context"
	"log"
	"net"
	"time"	
	"fmt"
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
	
	// client loop
	for {
		// input channel
		inputCh := make(chan string)
		go func() {
			var client_input string
			for {
				fmt.Scanln(&client_input)
				inputCh <- client_input
			}
		}()
	
		// write test message
		go func() {
			message := ("Hello from client!")
			
			if _, err := conn.Write([]byte(string(len(message)) + message)); err != nil {
				log.Fatal(err)
			}
		}()

		// quit client
		select {
		case input := <-inputCh:
			if input == "/quit" {
				fmt.Println("Quitting")
				return
			}
		}
	}
}
