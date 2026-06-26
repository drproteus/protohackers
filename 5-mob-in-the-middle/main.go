package main

import (
	"bufio"
	"log"
	"net"
	"regexp"
)

const UPSTREAM_ADDRESS = "chat.protohackers.com:16963"
const BOGUSCOIN_ADDRESS = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
const BOGUSCOIN_PATTERN = `\b7[a-zA-Z0-9]{25,34}\b`

func main() {
	listener, err := net.Listen("tcp", ":13337")
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
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return err
	} else {
		log.Println("Response relayed.")
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	pconn, err := net.Dial("tcp", UPSTREAM_ADDRESS)
	if err != nil {
		log.Printf("Error connecting to upstream chat server: %v\n", err)
		return
	} else {
		log.Println("Connected to chat server.")
	}
	defer pconn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Error reading from downstream client: %v\n", err)
			return
		}
		shakedown(&message)
		err = writeResponse(pconn, message)
		if err != nil {
			return
		}
	}
}

func shakedown(message *[]byte) {
	re, _ := regexp.Compile(BOGUSCOIN_PATTERN)
	re.ReplaceAll(*message, []byte(BOGUSCOIN_ADDRESS))
}
