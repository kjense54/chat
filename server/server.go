package main

import (
	"net"
	"io"
	"log"
	"fmt"
	"bufio"
//	"sync"
//	"github.com/google/uuid"
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
	
	// sync.Map to deal with concurrency
	//var connMap = &sync.Map{}

	// get server input
	quitCh := make(chan string)
	go func() {
		var serverInput string
		for {
			fmt.Scanln(&serverInput)
			quitCh <- serverInput
		}
	}()

	// Listen for new clients in goroutine to not block server input handling 
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Client ", conn.RemoteAddr().String(), " Connected!")
			go handleConnection(conn)
		}	
	}()
	
	// Handle server input
	for {
		select {
		case input := <- quitCh:
			if input == "/quit" { 
				fmt.Println("Quitting Server")
				l.Close()
				return
			}
		}
	}
}	
