package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"
	"net"
)

type Camera struct {
	ID    int
	Road  uint16 // in [0-65535]
	Mile  uint16 // mile marker on road
	Limit uint16 // speed limit
}

// Map of Camera IDs to Cameras
type State struct {
	Cameras map[int]Camera
}

type MessageType int

const (
	ErrorMessage MessageType = iota
	PlateMessage
	TicketMessage
	WantHeartbeat
	Heartbeat
	IAmCamera
	IAmDispatcher
)

var messageTypeName = map[MessageType]string{
	ErrorMessage:  "Error",
	PlateMessage:  "Plate",
	TicketMessage: "Ticket",
	WantHeartbeat: "WantHeartbeat",
	Heartbeat:     "Heartbeat",
	IAmCamera:     "IAmCamera",
	IAmDispatcher: "IAmDispatcher",
}

var messageTypeByte = map[MessageType]byte{
	ErrorMessage:  0x10,
	PlateMessage:  0x20,
	TicketMessage: 0x21,
	WantHeartbeat: 0x40,
	Heartbeat:     0x41,
	IAmCamera:     0x80,
	IAmDispatcher: 0x81,
}

func main() {
	listener, err := net.Listen("tcp", "6969")
	if err != nil {
		log.Fatal("Error creating server: ", err)
	}
	defer listener.Close()
	state := State{make(map[int]Camera)}
	cameraId := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handle(conn, &state, cameraId)
		cameraId++
	}
	// state := State{make(map[int]Camera)}
	// c1 := Camera{1, 20, 60}
	// c2 := Camera{1, 30, 50}
	// c3 := Camera{2, 10, 35}
	// c4 := Camera{1, 100, 70}
	// registerCamera(c1, &state)
	// registerCamera(c2, &state)
	// registerCamera(c3, &state)
	// registerCamera(c4, &state)
	// log.Println(getMonitoredRoads(&state))
}

func handle(conn net.Conn, state *State, cameraId int) {
	reader := bufio.NewReader(conn)
	for {
		msgTypeByte, err := reader.ReadByte()
		if err != nil {
			log.Printf("Error reading message type: %v", err)
			return
		}
		messageType, err := getMessageType(msgTypeByte)
		if err != nil {
			log.Println("Error getting message type!")
			return
		}
		if messageType == IAmCamera {
			roadBytes := make([]byte, 2)
			mileBytes := make([]byte, 2)
			limitBytes := make([]byte, 2)
			for range 2 {
				b, _ := reader.ReadByte()
				roadBytes = append(roadBytes, b)
			}
			for range 2 {
				b, _ := reader.ReadByte()
				mileBytes = append(mileBytes, b)
			}
			for range 2 {
				b, _ := reader.ReadByte()
				limitBytes = append(limitBytes, b)
			}
			c := Camera{
				cameraId,
				binary.BigEndian.Uint16(roadBytes),
				binary.BigEndian.Uint16(mileBytes),
				binary.BigEndian.Uint16(limitBytes),
			}
			registerCamera(c, state)
		}
	}
}

// func getMonitoredRoads(state *State) []uint16 {
// 	roads := make([]uint16, 0)
// 	for i := range (*state).Cameras {
// 		roads = append(roads, uint16(i))
// 	}
// 	return roads
// }

func registerCamera(camera Camera, state *State) {
	(*state).Cameras[int(camera.ID)] = camera
}

func handleMessage(conn net.Conn, state *State, msgTypeByte byte) {
	var messageType MessageType
	for mType, mByte := range messageTypeByte {
		if mByte == msgTypeByte {
			messageType = mType
			break
		}
	}
	if messageType < 0 {
		log.Printf("Invalid message type: %x", msgTypeByte)
		return
	}
}

func getMessageType(msgTypeByte byte) (MessageType, error) {
	for mType, mByte := range messageTypeByte {
		if mByte == msgTypeByte {
			return mType, nil
		}
	}
	return 0, errors.New("Unexpected message type")
}

func writeResponse(conn net.Conn, msg string) {
	response := []byte{byte(len(msg))}
	response = append(response, []byte(msg)...)
	conn.Write(response)
}

func writeError(conn net.Conn, msg string) {
	mByte := messageTypeByte[ErrorMessage]
	mBytes := []byte{mByte}
	conn.Write(mBytes)
	writeResponse(conn, msg)
}

func issueTicket(
	conn net.Conn,
	plate string,
	road uint16,
	mile1 uint16,
	timestamp1 uint32,
	mile2 uint16,
	timestamp2 uint32,
	speed uint16) {
	// Set the message type to 'Ticket'
	mByte := messageTypeByte[TicketMessage]
	mBytes := []byte{mByte}
	// Write the length of the plate string
	mBytes = append(mBytes, byte(len(plate)))
	// Write the plate string bytes
	mBytes = append(mBytes, []byte(plate)...)
	// Append the other fields in order
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, road)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, mile1)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, timestamp1)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, mile2)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, timestamp2)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, speed)
	conn.Write(mBytes)
}
