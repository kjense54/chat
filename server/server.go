package main

import (
	"net"
	"io"
	"log"
	"fmt"
	"bufio"
	"sync"
	"github.com/google/uuid"
)

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
	var connMap = &sync.Map{}

	// get server input
	inputCh := make(chan string)
	go func() {
		var serverInput string
		for {
			fmt.Scanln(&serverInput)
			inputCh <- serverInput
		}
	}()

	// Listen for new clients in goroutine to not block server input handling 
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Client at ", conn.RemoteAddr().String(), " connected!")
			
			// store in map
			id := uuid.New().String()
			connMap.Store(id, conn)

			go handleConnection(id, conn, connMap)
		}	
	}()
	
	// Handle server input
	for {
		select {
		case input := <- inputCh:
			if input == "/quit" { 
				fmt.Println("Quitting Server")
				l.Close()
				return
			}
		}
	}
}	

// handle each client's connection
func handleConnection(id string, conn net.Conn, connMap *sync.Map) {
	defer func() {
		fmt.Println("Client at ", conn.RemoteAddr().String(), " disconnected")
		conn.Close()
		connMap.Delete(id)
	}()

	c := bufio.NewReader(conn)
	
	// send received data to all clients 
	for {
		//read first byte for message length
		size, err := c.ReadByte()
		if err != nil {
			return	
		}
		
		buff := make([]byte, int(size)) // largest msg size

		// read full message to buff
		_, err1 := io.ReadFull(c, buff)
		if err1 != nil {
			return
		}
		// print to server terminal
		fmt.Println(string(buff))

		// send to all clients
		connMap.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				if _, err := conn.Write([]byte(string(len(buff)) + string(buff))); err != nil {
					fmt.Println("Error writing to connection")
				} 
			}
			return true
		})
	}
}
