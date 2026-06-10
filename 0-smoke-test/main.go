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

func writeResponse(conn net.Conn, message []byte) error {
	_, err := conn.Write(message)
	log.Printf("[RECEIVED] %s", message)
	if err != nil {
		log.Printf("Server response error: %v", err)
		return err
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if message != nil {
			writeResponse(conn, message)
		}
		if err != nil {
			return
		}
	}
}
