package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/big"
	"net"
	"strconv"
	"strings"
)

type testRequest struct {
	Method string `json:"method"`
	Number string `json:"number"`
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
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if message != nil {
			log.Printf("Received '%s'", strings.TrimSpace(string(message)))
		}
		if err != nil {
			log.Printf("Failed to read request: %v", err)
			writer.WriteString(`{"error": "Malformed request"}`)
			writer.WriteString("\n")
			writer.Flush()
			return
		}
		request := testRequest{}
		err = json.Unmarshal(message, &request)
		if err != nil {
			log.Printf("Malformed request: %v", err)
			writer.WriteString(`{"error": "Malformed request"}`)
			writer.WriteString("\n")
			writer.Flush()
			return
		}
		if request.Method != "isPrime" {
			log.Printf("Invalid method: %s", request.Method)
			writer.WriteString(`{"error": "Malformed request"}`)
			writer.WriteString("\n")
			writer.Flush()
			return
		}
		num, err := strconv.ParseInt(request.Number, 10, 32)
		if err != nil {
			log.Printf("Invalid or missing number: %s", request.Number)
			writer.WriteString(`{"error": "Malformed request"}`)
			writer.WriteString("\n")
			writer.Flush()
			return
		}
		prime := big.NewInt(int64(num)).ProbablyPrime(0)
		response, err := json.Marshal(testResponse{Method: "isPrime", Prime: prime})
		log.Printf("Result: '%s'", response)
		writer.WriteString(string(response))
		writer.WriteString("\n")
		writer.Flush()
	}
}
