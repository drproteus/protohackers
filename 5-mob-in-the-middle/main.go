package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"regexp"
)

const UPSTREAM_ADDRESS = "chat.protohackers.com:16963"
const BOGUSCOIN_ADDRESS = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
const BOGUSCOIN_PATTERN = `(^|\s)(7[A-Za-z0-9]{25,34})`

func main() {
	listener, err := net.Listen("tcp", ":13337")
	if err != nil {
		log.Fatal("Error creating server: ", err)
	}
	defer listener.Close()
	connId := 0
	for {
		pconn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(pconn, connId)
		connId++
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

func handleConnection(pconn net.Conn, connId int) {
	defer pconn.Close()
	uconn, err := net.Dial("tcp", UPSTREAM_ADDRESS)
	if err != nil {
		log.Printf("Error connecting to upstream chat server: %v\n", err)
		return
	}
	go readLoop(pconn, uconn, connId)
	go writeLoop(pconn, uconn, connId)
	for {
	}
}

func readLoop(pconn net.Conn, uconn net.Conn, connId int) {
	defer uconn.Close()
	reader := bufio.NewReader(pconn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("[%d] (u) <-- (p) %v\n", connId, err)
			return
		}
		log.Printf("[%d] (p) --> (u) %s", connId, message)
		err = writeResponse(uconn, shakedown(message))
		if err != nil {
			log.Printf("[%d] (p) --> (u) %v\n", connId, err)
			return
		}
	}
}

func writeLoop(pconn net.Conn, uconn net.Conn, connId int) {
	reader := bufio.NewReader(uconn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("[%d] (p) <-- (u) %v\n", connId, err)
			return
		}
		log.Printf("[%d] (u) --> (p) %s", connId, message)
		err = writeResponse(pconn, shakedown(message))
		if err != nil {
			log.Printf("[%d] (u) --> (p) %v\n", connId, err)
			return
		}
	}
}

func shakedown(message []byte) []byte {
	re := regexp.MustCompile(`(^|\s)(7[A-Za-z0-9]{25,34})`)
	b := new(bytes.Buffer)
	last := 0
	for _, loc := range re.FindAllSubmatchIndex(message, -1) {
		tokenEnd := loc[5]
		// Verify right boundary: byte after token must be non-alphanumeric or EOS
		if tokenEnd < len(message) && !isSpace(message[tokenEnd]) {
			continue // token is part of a longer alphanumeric run — skip
		}
		b.Write(message[last:loc[4]]) // includes left boundary
		b.WriteString(BOGUSCOIN_ADDRESS)
		last = tokenEnd
	}
	b.Write(message[last:])
	return b.Bytes()
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}
