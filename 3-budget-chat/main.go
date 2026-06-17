package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"unicode"
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
		// conn.SetDeadline(time.Now().Add(5 * time.Second))
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
		log.Print(msg)
		rw.WriteString(msg)
		rw.Flush()
		conn.Close()
		return
	}
	if containsNonAlphanumeric(name) {
		msg := "Illegal name!\n"
		log.Print(msg)
		rw.WriteString(msg)
		rw.Flush()
		conn.Close()
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
	users := getConnectedUsers(state, connId)
	usersString := strings.Join(users, ", ")
	rw.WriteString(fmt.Sprintf("* The room contains: %s\n", usersString))
	rw.Flush()
	go readLoop(state, connId)
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
			message := msg.message
			if msg.name != "_SYSTEM" {
				message = fmt.Sprintf("[%s] %s", msg.name, msg.message)
			}
			_, err := state.rwMap[i].WriteString(fmt.Sprintf("%s\n", message))
			if err != nil {
				// connection is probably closed (TODO)
				continue
			}
			err = state.rwMap[i].Flush()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Wrote to conn ID %d, %s\n", i, state.nameMap[i])
		}
		time.Sleep(WRITE_LOOP_SLEEP * time.Second)
	}
}

func cleanUpConn(state *ChatState, connId int) {
	state.rwMap[connId].Flush()
	state.connMap[connId].Close()
	msg := ChatMessage{connId, "_SYSTEM", fmt.Sprintf("* %s has left the room", state.nameMap[connId])}
	state.msgQueue = append(state.msgQueue, msg)
	// just clear the name in the map to indicate disconnect (TODO, improve)
	state.nameMap[connId] = ""
}

func readLoop(state *ChatState, connId int) {
	defer cleanUpConn(state, connId)
	rw := state.rwMap[connId]
	for {
		message, err := rw.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		message = strings.TrimSpace(message)
		state.msgQueue = append(state.msgQueue, ChatMessage{connId, state.nameMap[connId], message})
		time.Sleep(READ_LOOP_SLEEP * time.Second)
	}
}

func getConnectedUsers(state *ChatState, connId int) []string {
	users := make([]string, 0)
	for i := 0; i < state.connId; i++ {
		if i == connId || state.nameMap[connId] == "" {
			continue
		}
		users = append(users, state.nameMap[i])
	}
	return users
}

func containsNonAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true // Found a non-alphanumeric character
		}
	}
	return false
}
