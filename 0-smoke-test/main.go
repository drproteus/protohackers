package main

import (
	"bufio"
	"io"
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

func writeResponse(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Printf("Server response error: %v", err)
		return err
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	data := ""
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Println("Connection terminated by EOF")
			writeResponse(conn, data)
			return
		}
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}
		data += message
		writeResponse(conn, message)
	}
}
