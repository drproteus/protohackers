package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"
	"net"
)

type ConnID int

type ClientType int

const (
	CameraClient ClientType = iota
	DispatcherClient
)

type Client struct {
	ID                ConnID
	Type              ClientType
	HeartbeatInterval uint32 // heartbeat interval, if set
	Connected         bool   // if connected
}

type Observation struct {
	Plate     string
	Timestamp uint32
}

type Camera struct {
	Road    uint16 // in [0-65535]
	Mile    uint16 // mile marker on road
	Limit   uint16 // speed limit
	Enabled bool   // if the camera is valid and was registered
}

type Dispatcher struct {
	Roads   []uint16
	Enabled bool
}

// Map of Camera IDs to Cameras
type State struct {
	Cameras      map[ConnID]Camera        // map of connID to cameras
	Dispatchers  map[ConnID]Dispatcher    // map of connID of dispatchers
	Observations map[uint16][]Observation // map of road ID to plate/timestamps on that road
	Clients      map[ConnID]Client        // map of connID to client
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
	state := State{
		make(map[ConnID]Camera),
		make(map[ConnID]Dispatcher),
		make(map[uint16][]Observation),
		make(map[ConnID]Client),
	}
	connId := ConnID(0)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handle(conn, &state, connId)
		connId++
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

func handle(conn net.Conn, state *State, connId ConnID) {
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
		switch messageType {
		case IAmCamera:
			handleIAmCamera(*reader, state, connId)
		case IAmDispatcher:
			handleIAmDispatcher(*reader, state, connId)
		case WantHeartbeat:
			handleWantHeartbeat(*reader, state, connId)
		case PlateMessage:
			handlePlateMessage(*reader, state, connId)
		}
	}
}

// func getMonitoredRoads(state *State) []int {
// 	roads := make([]int, 0)
// 	for i := range (*state).Cameras {
// 		if slices.Contains(roads, i) {
// 			continue
// 		}
// 		roads = append(roads, i)
// 	}
// 	return roads
// }

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

func sendHeartbeat(conn net.Conn) {
	mBytes := []byte{messageTypeByte[Heartbeat]}
	conn.Write(mBytes)
}

func handleIAmCamera(reader bufio.Reader, state *State, connId ConnID) {
	c := (*state).Cameras[connId]
	if c.Enabled {
		// ERROR, camera already registered
		return
	}
	d := (*state).Dispatchers[connId]
	if d.Enabled {
		// Error, connection already registered as dispatcher
		return
	}
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
	c = Camera{
		binary.BigEndian.Uint16(roadBytes),
		binary.BigEndian.Uint16(mileBytes),
		binary.BigEndian.Uint16(limitBytes),
		true,
	}
	(*state).Cameras[connId] = c
	client := Client{connId, CameraClient, 0, true}
	(*state).Clients[connId] = client
}

func handleIAmDispatcher(reader bufio.Reader, state *State, connId ConnID) {
	d := (*state).Dispatchers[connId]
	if d.Enabled {
		// ERROR, dispatcher already registered
		return
	}
	c := (*state).Cameras[connId]
	if c.Enabled {
		// Error, connection already registered as camera
		return
	}
	numRoadByte, _ := reader.ReadByte()
	numRoads := int(numRoadByte)
	roads := make([]uint16, numRoads)
	for range numRoads {
		roadBytes := make([]byte, 2)
		for range 2 {
			b, _ := reader.ReadByte()
			roadBytes = append(roadBytes, b)
		}
		road := binary.BigEndian.Uint16(roadBytes)
		roads = append(roads, road)
	}
	d = Dispatcher{roads, true}
	(*state).Dispatchers[connId] = d
	client := Client{connId, DispatcherClient, 0, true}
	(*state).Clients[connId] = client
}

func handleWantHeartbeat(reader bufio.Reader, state *State, connId ConnID) {
	cl := (*state).Clients[connId]
	if !cl.Connected {
		// ERROR, requested heartbeat on unregistered conn ID
		return
	}
	intervalBytes := make([]byte, 4)
	for range 4 {
		b, _ := reader.ReadByte()
		intervalBytes = append(intervalBytes, b)
	}
	interval := binary.BigEndian.Uint32(intervalBytes)
	cl.HeartbeatInterval = interval
	(*state).Clients[connId] = cl
}

func handlePlateMessage(reader bufio.Reader, state *State, connId ConnID) {
	camera := (*state).Cameras[connId]
	if !camera.Enabled {
		// ERROR, camera is not connected
		return
	}
	plateLengthByte, _ := reader.ReadByte()
	plateLength := int(plateLengthByte)
	plateBytes := make([]byte, plateLength)
	for range plateLength {
		b, _ := reader.ReadByte()
		plateBytes = append(plateBytes, b)
	}
	timestampBytes := make([]byte, 4)
	for range 4 {
		b, _ := reader.ReadByte()
		timestampBytes = append(timestampBytes, b)
	}
	timestamp := binary.BigEndian.Uint32(timestampBytes)
	ro := (*state).Observations[camera.Road]
	ro = append(ro, Observation{string(plateBytes), timestamp})
	(*state).Observations[camera.Road] = ro
}
