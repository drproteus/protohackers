package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const WRITE_LOOP_SLEEP = 1 // seconds
const READ_LOOP_SLEEP = 1  // seconds

type ChatMessage struct {
	connId  int
	name    string
	message string
}

type ChatState struct {
	connId   int
	nameMap  map[int]string
	connMap  map[int]net.Conn
	rwMap    map[int]*bufio.ReadWriter
	history  []ChatMessage
	msgQueue []ChatMessage
}

func main() {
	state := ChatState{connId: 0}
	state.nameMap = make(map[int]string)
	state.connMap = make(map[int]net.Conn)
	state.rwMap = make(map[int]*bufio.ReadWriter)
	state.history = []ChatMessage{}
	state.msgQueue = make([]ChatMessage, 0)
	listener, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Fatal("Error creating server: ", err)
	}
	defer listener.Close()
	go writeLoop(&state)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		state.connMap[state.connId] = conn
		writer := bufio.NewWriter(conn)
		reader := bufio.NewReader(conn)
		rw := bufio.NewReadWriter(reader, writer)
		state.rwMap[state.connId] = rw
		go handleConnection(conn, state.connId, &state)
		state.connId += 1
	}
}

func handleConnection(conn net.Conn, connId int, state *ChatState) {
	defer conn.Close()
	rw := state.rwMap[connId]
	rw.WriteString("Welcome to budgetchat! What shall I call you?\n")
	rw.Flush()
	name, err := rw.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	name = strings.TrimSpace(name)
	if len(name) < 1 {
		msg := "Name too short!\n"
		log.Println(msg)
		rw.WriteString(msg)
		rw.Flush()
		return
	}
	log.Printf("[USER JOINED] [CONN ID: %d] [NAME: %s]", connId, name)
	state.nameMap[connId] = name
	msg := ChatMessage{connId, "_SYSTEM", fmt.Sprintf("* %s has entered the room", name)}
	state.msgQueue = append(state.msgQueue, msg)
	// alert other users of new user
	// for i := 0; i <= state.connId; i++ {
	// 	if i != connId && state.nameMap[connId] != "" {
	// 		handleUpdate(state, i, fmt.Sprintf("* %s has entered the room", name))
	// 	}
	// }
	// for {
	// }
	go readLoop(*state, state.connId)
}

func writeLoop(state *ChatState) {
	for {
		if len(state.msgQueue) < 1 {
			continue
		}
		// get first queued message
		msg := state.msgQueue[0]
		log.Printf("[%s][%d] %s\n", msg.name, msg.connId, msg.message)
		// remove it from the queue
		state.msgQueue = state.msgQueue[1:]
		// loop through active connections to broadcast,
		// skipping the connection that sent it
		for i := 0; i <= state.connId; i++ {
			if state.nameMap[i] == "" || msg.connId == i {
				continue
			}
			state.rwMap[i].WriteString(fmt.Sprintf("%s\n", msg.message))
			state.rwMap[i].Flush()
			log.Printf("Wrote to conn ID %d, %s\n", i, state.nameMap[i])
		}
		time.Sleep(WRITE_LOOP_SLEEP * time.Second)
	}
}

func readLoop(state ChatState, connId int) {
	for {
		time.Sleep(READ_LOOP_SLEEP * time.Second)
	}
}
