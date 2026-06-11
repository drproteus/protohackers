package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/big"
	"net"
	"strings"
)

type testRequest struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

type testResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	listener, err := net.Listen("tcp", ":8000")
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
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		message, err := rw.ReadBytes('\n')
		if message != nil {
			log.Printf("Received '%s'", strings.TrimSpace(string(message)))
		}
		if err != nil {
			log.Printf("Failed to read request: %v", err)
			rw.WriteString(`{"error": "Malformed request"}`)
			return
		}
		request := testRequest{}
		err = json.Unmarshal(message, &request)
		if err != nil {
			log.Printf("Malformed request: %v", err)
			rw.WriteString(`{"error": "Malformed request"}`)
			return
		}
		if request.Method != "isPrime" {
			log.Printf("Invalid method: %s", request.Method)
			rw.WriteString(`{"error": "Malformed request"}`)
			return
		}
		prime := big.NewInt(int64(request.Number)).ProbablyPrime(0)
		response, err := json.Marshal(testResponse{Method: "isPrime", Prime: prime})
		log.Printf("Result: '%s'", response)
		_, err = rw.WriteString(string(response))
		if err != nil {
			log.Panic(err)
			return
		}
	}
}
