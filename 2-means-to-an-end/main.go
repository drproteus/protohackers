package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
)

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

func writeResponse(conn net.Conn, mean int32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(mean))
	_, err := conn.Write(buf)
	if err != nil {
		log.Printf("Failure to communicate: %v", err)
		return err
	}
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	record := make(map[int32]int32)
	for {
		msgTypeBytes := make([]byte, 1)
		_, err := io.ReadFull(reader, msgTypeBytes)
		if err != nil {
			log.Printf("Error reading first byte of message: %v\n", err)
			return
		}
		msgType := string(msgTypeBytes)
		if msgType == "I" {
			// The first byte is equivalent to ASCII 'I',
			// This is an insertion.
			timestampBytes := make([]byte, 4)
			_, err = io.ReadFull(reader, timestampBytes)
			timestampUnsigned := binary.BigEndian.Uint32(timestampBytes)
			timestamp := int32(timestampUnsigned)
			priceBytes := make([]byte, 4)
			_, err = io.ReadFull(reader, priceBytes)
			priceUnsigned := binary.BigEndian.Uint32(priceBytes)
			price := int32(priceUnsigned)
			record[timestamp] = price
		} else if msgType == "Q" {
			// The first byte is equivalent to ASCII 'Q',
			// This is a query.
			mintimeBytes := make([]byte, 4)
			_, err = io.ReadFull(reader, mintimeBytes)
			mintimeUnsigned := binary.BigEndian.Uint32(mintimeBytes)
			mintime := int32(mintimeUnsigned)
			maxtimeBytes := make([]byte, 4)
			_, err = io.ReadFull(reader, maxtimeBytes)
			maxtimeUnsigned := binary.BigEndian.Uint32(maxtimeBytes)
			maxtime := int32(maxtimeUnsigned)
			if mintime > maxtime {
				writeResponse(conn, 0)
				continue
			}
			// Now it's time to crunch the numbers.
			var sum int64
			var count int64
			for timestamp, price := range record {
				if timestamp < mintime || timestamp > maxtime {
					continue
				}
				count++
				sum += int64(price)
			}
			if count < 1 {
				writeResponse(conn, 0)
				continue
			}
			mean := int32(sum / count)
			writeResponse(conn, mean)
		} else {
			log.Printf("Unexpected message type: %v", msgType)
		}
	}
}
