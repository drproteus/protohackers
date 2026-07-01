package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"slices"
	"time"
)

type ConnID int

type ClientType int

type Client struct {
	ID                ConnID
	HeartbeatInterval uint32 // heartbeat interval, if set
	Connected         bool   // if connected
	Camera            *Camera
	Dispatcher        *Dispatcher
	Conn              *net.Conn
}

type Observation struct {
	Road      uint16
	Mile      uint16
	Plate     string
	Timestamp uint32
}

type Camera struct {
	Road    uint16 // in [0-65535]
	Mile    uint16 // mile marker on road
	Limit   uint16 // speed limit
	Enabled bool   // if the camera is valid and was registered
	Conn    *net.Conn
}

type Dispatcher struct {
	Roads   []uint16
	Enabled bool
	Conn    *net.Conn
}

type State struct {
	Cameras      []Camera
	Dispatchers  []Dispatcher
	Observations []Observation
	Clients      map[ConnID]Client // map of connID to client
	TicketQ      []Ticket          // ticket queue
}

type Ticket struct {
	Plate      string
	Road       uint16
	Mile1      uint16
	Timestamp1 uint32
	Mile2      uint16
	Timestamp2 uint32
	Speed      uint16
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
	listener, err := net.Listen("tcp", ":13337")
	if err != nil {
		log.Fatal("Error creating server: ", err)
	}
	defer listener.Close()
	state := State{
		make([]Camera, 0),
		make([]Dispatcher, 0),
		make([]Observation, 0),
		make(map[ConnID]Client),
		make([]Ticket, 0),
	}
	go issueTickets(&state)
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
			log.Printf("Error reading message type: %v\n", err)
			return
		}
		messageType, err := getMessageType(msgTypeByte)
		if err != nil {
			log.Printf("Error getting message type: %x\n", msgTypeByte)
			return
		}
		switch messageType {
		case IAmCamera:
			log.Println("New camera")
			handleIAmCamera(conn, *reader, state, connId)
		case IAmDispatcher:
			log.Println("New dispatcher")
			handleIAmDispatcher(conn, *reader, state, connId)
		case WantHeartbeat:
			log.Println("Heartbeat request")
			handleWantHeartbeat(conn, *reader, state, connId)
		case PlateMessage:
			log.Println("Plate message")
			handlePlateMessage(conn, *reader, state, connId)
		}
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

func issueTicket(conn net.Conn, ticket Ticket) {
	// Set the message type to 'Ticket'
	mByte := messageTypeByte[TicketMessage]
	mBytes := []byte{mByte}
	// Write the length of the plate string
	mBytes = append(mBytes, byte(len(ticket.Plate)))
	// Write the plate string bytes
	mBytes = append(mBytes, []byte(ticket.Plate)...)
	// Append the other fields in order
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, ticket.Road)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, ticket.Mile1)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, ticket.Timestamp1)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, ticket.Mile2)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, ticket.Timestamp2)
	mBytes, _ = binary.Append(mBytes, binary.BigEndian, ticket.Speed)
	conn.Write(mBytes)
}

func sendHeartbeat(conn net.Conn) error {
	mBytes := []byte{messageTypeByte[Heartbeat]}
	_, err := conn.Write(mBytes)
	return err
}

func handleIAmCamera(conn net.Conn, reader bufio.Reader, state *State, connId ConnID) {
	cl := (*state).Clients[connId]
	if cl.Connected {
		// Error, connection already registered
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
	c := Camera{
		binary.BigEndian.Uint16(roadBytes),
		binary.BigEndian.Uint16(mileBytes),
		binary.BigEndian.Uint16(limitBytes),
		true,
		&conn,
	}
	client := Client{
		connId,
		0,
		true,
		&c,
		nil,
		&conn,
	}
	(*state).Clients[connId] = client
	(*state).Cameras = append((*state).Cameras, c)
}

func handleIAmDispatcher(conn net.Conn, reader bufio.Reader, state *State, connId ConnID) {
	cl := (*state).Clients[connId]
	if cl.Connected {
		// Error, connection already registered
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
	d := Dispatcher{roads, true, &conn}
	client := Client{
		connId,
		0,
		true,
		nil,
		&d,
		&conn,
	}
	(*state).Dispatchers = append((*state).Dispatchers, d)
	(*state).Clients[connId] = client
}

func handleWantHeartbeat(conn net.Conn, reader bufio.Reader, state *State, connId ConnID) {
	intervalBytes := make([]byte, 4)
	for range 4 {
		b, _ := reader.ReadByte()
		intervalBytes = append(intervalBytes, b)
	}
	interval := binary.BigEndian.Uint32(intervalBytes)
	cl := (*state).Clients[connId]
	if cl.Connected && cl.HeartbeatInterval > 0 {
		// Error, already registered for heartbeat
		return
	}
	if !cl.Connected {
		cl = Client{connId, interval, true, nil, nil, &conn}
		(*state).Clients[connId] = cl
	}
	cl.HeartbeatInterval = interval
	startHeartbeat(conn, state, connId, interval)
}

func startHeartbeat(conn net.Conn, state *State, connId ConnID, interval uint32) {
	if interval <= 0 {
		return
	}
	time.AfterFunc(time.Duration(interval)*time.Second/10, func() {
		cl := (*state).Clients[connId]
		if !cl.Connected {
			return
		}
		err := sendHeartbeat(conn)
		if err != nil {
			startHeartbeat(conn, state, connId, interval)
		}
	})
}

func handlePlateMessage(conn net.Conn, reader bufio.Reader, state *State, connId ConnID) {
	cl := (*state).Clients[connId]
	if !cl.Connected {
		// ERROR, client not connected
	}
	camera := cl.Camera
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
	ob := Observation{
		camera.Road,
		camera.Mile,
		string(plateBytes),
		timestamp,
	}
	(*state).Observations = append((*state).Observations, ob)
	checkSpeed(state, ob)
}

func checkSpeed(state *State, ob Observation) {
	for _, ob2 := range (*state).Observations {
		if ob2.Plate != ob.Plate {
			continue
		}
		if ob2.Mile == ob.Mile {
			continue
		}
		for _, camera := range (*state).Cameras {
			maxSpeed := uint16(0)
			if camera.Mile > ob.Mile {
				continue
			}
			totalDistance := ob.Mile - camera.Mile
			totalTime := ob.Timestamp - ob2.Timestamp
			speed := uint16(int(totalDistance) / int(totalTime))
			if speed > maxSpeed {
				maxSpeed = speed
			}
			if maxSpeed >= camera.Limit {
				// You've violated the law!
				t := Ticket{
					ob.Plate,
					ob.Road,
					ob.Mile,
					ob.Timestamp,
					ob2.Mile,
					ob2.Timestamp,
					maxSpeed,
				}
				log.Printf("Traffic violation: %s\n", ob.Plate)
				(*state).TicketQ = append((*state).TicketQ, t)
			}
		}
	}
}

func issueTickets(state *State) {
	for {
		if len((*state).TicketQ) < 1 {
			continue
		}
		t := (*state).TicketQ[0]
		(*state).TicketQ = (*state).TicketQ[1:]
		didIssue := false
		for _, d := range (*state).Dispatchers {
			if !slices.Contains(d.Roads, t.Road) {
				continue
			}
			issueTicket(*d.Conn, t)
			didIssue = true
		}
		// Put it back if no dispatcher available
		if !didIssue {
			(*state).TicketQ = append((*state).TicketQ, t)
		}
	}
}
