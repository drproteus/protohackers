package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var lock = sync.Mutex{}

func main() {
	pc, err := net.ListenPacket("udp", ":13337")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	store := map[string]string{"version": "Ken's Key-Value Store 1.0"}

	for {
		buf := make([]byte, 1000)
		nBytes, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Println("read:", err)
			continue
		}
		go respond(pc, addr, buf[:nBytes], store)
	}
}

func respond(pc net.PacketConn, addr net.Addr, payload []byte, store map[string]string) {
	// msg := strings.TrimSpace(string(payload))
	// msg := strings.TrimSuffix(string(payload), "\n")
	msg := string(payload)
	log.Println("msg:", msg)
	sepIndex := strings.IndexRune(msg, '=')
	log.Println("sepIndex:", sepIndex)
	if sepIndex < 0 {
		read(pc, addr, msg, store)
	} else {
		key := msg[:sepIndex]
		if key != "version" {
			value := msg[sepIndex+1:]
			write(pc, addr, key, value, store)
			log.Printf("Value of %s is now %s\n", key, value)
		}
	}
}

func read(pc net.PacketConn, addr net.Addr, key string, store map[string]string) {
	lock.Lock()
	defer lock.Unlock()
	value := store[key]
	log.Printf("Value of %s is %s\n", key, value)
	_, err := pc.WriteTo([]byte(fmt.Sprintf("%s=%s", key, value)), addr)
	if err != nil {
		log.Println("write:", err)
	}
}

func write(pc net.PacketConn, addr net.Addr, key string, value string, store map[string]string) {
	lock.Lock()
	defer lock.Unlock()
	store[key] = value
}
