package main

import (
	"fmt"
	"net"
	"log"
	//"context"
)

func main() {
	fmt.Println("Client Started")

	localAddr := "localhost:8090"
	conn, err := net.Dial("tcp", localAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("Hello!")); err != nil {
		log.Fatal(err)
	}
}
