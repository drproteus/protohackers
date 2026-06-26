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

type ProxyMessage struct {
	Original string
	Modified string
}

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
	// reader := bufio.NewReader(conn)
	// for {
	// 	message, err := reader.ReadBytes('\n')
	// 	if err != nil {
	// 		log.Printf("Error reading from downstream client: %v\n", err)
	// 		return
	// 	}
	// 	shakedown(&message)
	// 	err = writeResponse(pconn, message)
	// 	if err != nil {
	// 		return
	// 	}
	// }
	msgQueue := make([]ProxyMessage, 0)
	go readLoop(pconn, &msgQueue)
	go writeLoop(conn, pconn, &msgQueue)
	for {
		// Keep alive?
	}
}

func readLoop(pconn net.Conn, msgQueue *[]ProxyMessage) {
	reader := bufio.NewReader(pconn)
	for {
		message, _ := reader.ReadBytes('\n')
		log.Printf("Received: %s\n", message)
		pMsg := ProxyMessage{}
		pMsg.Original = string(message)
		pMsg.Modified = string(shakedown(message))
		log.Printf("Substituted message: %s\n", message)
		*msgQueue = append(*msgQueue, pMsg)
	}
}

func writeLoop(conn net.Conn, pconn net.Conn, msgQueue *[]ProxyMessage) {
	writer := bufio.NewWriter(conn)
	pwriter := bufio.NewWriter(pconn)
	for {
		if len(*msgQueue) < 1 {
			continue
		}
		msg := (*msgQueue)[0]
		*msgQueue = (*msgQueue)[1:]
		writer.WriteString(msg.Original)
		writer.Flush()
		pwriter.WriteString(msg.Modified)
		pwriter.Flush()
		log.Printf("Wrote: %s\n", msg)
	}
}

func shakedown(message []byte) []byte {
	re, _ := regexp.Compile(BOGUSCOIN_PATTERN)
	re.ReplaceAll(message, []byte(BOGUSCOIN_ADDRESS))
	return message
}
