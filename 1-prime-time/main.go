package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/big"
	"net"
	"strings"

	"github.com/oapi-codegen/nullable"
)

type testRequest struct {
	Method string                 `json:"method"`
	Number nullable.Nullable[int] `json:"number"`
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

func writeOut(writer bufio.Writer, message string) {
	writer.WriteString(message)
	writer.WriteString("\n")
	writer.Flush()
}

func errorOut(writer bufio.Writer) {
	writeOut(writer, `{"error": "Malformed request"}`)
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
			errorOut(*writer)
			return
		}
		request := testRequest{}
		err = json.Unmarshal(message, &request)
		if err != nil || request.Method != "isPrime" || !request.Number.IsSpecified() {
			errorOut(*writer)
			return
		}
		num, _ := request.Number.Get()
		prime := big.NewInt(int64(num)).ProbablyPrime(0)
		response, err := json.Marshal(testResponse{Method: "isPrime", Prime: prime})
		log.Printf("Result: '%s'", response)
		writeOut(*writer, string(response))
	}
}
