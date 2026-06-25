package main

import (
	"bufio"
	"encoding/binary"
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
		msgQueryType, err := reader.ReadByte()
		if err != nil {
			log.Printf("Error reading first byte of message: %v\n", err)
			return
		}
		if int(msgQueryType) == 73 {
			// The first byte is equivalent to ASCII 'I',
			// This is an insertion.
			timestampBytes := make([]byte, 4)
			for i := 0; i < 4; i++ {
				nextByte, err := reader.ReadByte()
				if err != nil {
					log.Printf("Error reading timestamp: %v\n", err)
					return
				}
				timestampBytes = append(timestampBytes, nextByte)
			}
			timestampUnsigned := binary.BigEndian.Uint32(timestampBytes)
			timestamp := int32(timestampUnsigned)
			priceBytes := make([]byte, 4)
			for i := 0; i < 4; i++ {
				nextByte, err := reader.ReadByte()
				if err != nil {
					log.Printf("Error reading price: %v\n", err)
					return
				}
				priceBytes = append(priceBytes, nextByte)
			}
			priceUnsigned := binary.BigEndian.Uint32(priceBytes)
			price := int32(priceUnsigned)
			record[timestamp] = price
		} else if int(msgQueryType) == 81 {
			// The first byte is equivalent to ASCII 'Q',
			// This is a query.
			mintimeBytes := make([]byte, 4)
			for i := 0; i < 4; i++ {
				nextByte, err := reader.ReadByte()
				if err != nil {
					log.Printf("Error reading mintime: %v\n", err)
					return
				}
				mintimeBytes = append(mintimeBytes, nextByte)
			}
			mintimeUnsigned := binary.BigEndian.Uint32(mintimeBytes)
			mintime := int32(mintimeUnsigned)
			maxtimeBytes := make([]byte, 4)
			for i := 0; i < 4; i++ {
				nextByte, err := reader.ReadByte()
				if err != nil {
					log.Printf("Error reading mintime: %v\n", err)
					return
				}
				maxtimeBytes = append(maxtimeBytes, nextByte)
			}
			maxtimeUnsigned := binary.BigEndian.Uint32(maxtimeBytes)
			maxtime := int32(maxtimeUnsigned)
			if mintime > maxtime {
				writeResponse(conn, 0)
				continue
			}
			// Now it's time to crunch the numbers.
			var sum int32
			var count int32
			for timestamp, price := range record {
				if timestamp < mintime || timestamp > maxtime {
					continue
				}
				count++
				sum += price
			}
			if count < 1 {
				writeResponse(conn, 0)
				continue
			}
			mean := sum / count
			writeResponse(conn, mean)
		} else {
			log.Printf("Unexpected message type, closing connection!")
		}
	}
}
