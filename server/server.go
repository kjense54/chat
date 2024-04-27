package main

import (
	"net"
	"io"
	"log"
	"fmt"
	"bufio"
)

func handleConnection(conn net.Conn) {
	defer func() {
		fmt.Println("Client at ", conn.RemoteAddr().String(), " disconnected")
		conn.Close()
	}()

	buff := make([]byte, 50) // largest msg
	c := bufio.NewReader(conn)
	
	// receive loop
	for {
		//read first byte for message length
		size, err := c.ReadByte()
		if err != nil {
			fmt.Println("Client disconnected")
			return	
		}
		// read full message
		_, err1 := io.ReadFull(c, buff[:int(size)])
		if err1 != nil {
			return
		}
		fmt.Println(string(buff[:int(size)]))
	}
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


	for {
		fmt.Println("Server Started")
		// input channel
		inputCh := make(chan string)
		go func() {
			var server_input string 
			for {
				fmt.Scanln(&server_input)
				inputCh <- server_input
			}
		}()
		
		// wait for connection
		go func() {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Client ", conn.RemoteAddr().String(), " Connected!")
			go handleConnection(conn)
		}()

		//quit server
		select {
		case input := <-inputCh:
			if input == "/quit" {
				fmt.Println("Quitting")
				return
			}
		}

	}
}	
