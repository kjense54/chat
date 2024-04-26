package main

import (
	"net"
	"io"
	"log"
	"fmt"
)

func handleConnection(c net.Conn) {
	io.Copy(c, c) // echo incoming data
}

func handleListen(l net.Listener) {
	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Client connected!")
	go handleConnection(conn)
}

func main() {

	// listen on tcp port 8090	
	localAddr := "localhost:8090"
	l, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatal(err)
	}
	// close listener after main ends
	defer l.Close()

	// receive server terminal input to channel
	inputCh := make(chan string) 
	go func() {
		var server_input string 
		for {
			fmt.Scanln(&server_input)
			inputCh <- server_input
		}
	}()

	for {
		go handleListen(l)

		//quit server
		select {
		case input := <-inputCh:
			if input == "quit" {
				fmt.Println("Quitting")
				return
			}
		}

	}
}	
