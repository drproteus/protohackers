package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatal("Error creating server: ", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Printf("Server response error: %v", err)
	}
}
