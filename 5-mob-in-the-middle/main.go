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
		pconn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(pconn)
	}
}

func writeResponse(conn net.Conn, message []byte) error {
	_, err := conn.Write(message)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return err
	}
	return nil
}

func handleConnection(pconn net.Conn) {
	defer pconn.Close()
	uconn, err := net.Dial("tcp", UPSTREAM_ADDRESS)
	if err != nil {
		log.Printf("Error connecting to upstream chat server: %v\n", err)
		return
	} else {
		log.Println("Connected to chat server.")
	}
	defer uconn.Close()

	go readLoop(pconn, uconn)
	go writeLoop(pconn, uconn)

	for {
	}
}

func readLoop(pconn net.Conn, uconn net.Conn) {
	reader := bufio.NewReader(pconn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("(p) %v\n", err)
			return
		}
		log.Printf("(p) --> %s", message)
		err = writeResponse(uconn, shakedown(message))
		if err != nil {
			log.Printf("(u) %v\n", err)
		}
	}
}

func writeLoop(pconn net.Conn, uconn net.Conn) {
	reader := bufio.NewReader(uconn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("(u) %v\n", err)
			return
		}
		log.Printf("(u) --> %s", message)
		err = writeResponse(pconn, shakedown(message))
		if err != nil {
			log.Printf("(p) %v\n", err)
			return
		}
	}
}

func shakedown(message []byte) []byte {
	re, _ := regexp.Compile(BOGUSCOIN_PATTERN)
	return re.ReplaceAll(message, []byte(BOGUSCOIN_ADDRESS))
}
