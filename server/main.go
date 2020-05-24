/*
server:
    maintain master state
    always be listening
    broadcast on a scheduled

    threads:
        listener:
            adds payloads to queue
        input processor:
            reads from input queue. Updates global state
        publisher:
            60x per second publish state to all clients
    state:
        each client/tank:
            id int //assigned at connection time
            address
            max frame processed
            last heard from
            tank data

*/

//nc -u localhost 1053
// https://ops.tips/blog/udp-client-and-server-in-go/
//https://gist.github.com/miekg/d9bc045c89578f3cc66a214488e68227
package main

import (

	// "encoding/binary"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/willhickey/TankWar/state"
)

const (
	clientTimeoutMillis = 5000
)

var (
	clients              = make(map[uint8]state.ClientState)
	listenQueue          = make(chan *state.Tank, 10)
	handshakeQueue       = make(chan *net.Addr, 10)
	nextClientId   int32 = 0
)

func main() {
	// listen to incoming udp packets
	conn, err := net.ListenPacket("udp", ":1984")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go dropDisconnectedClients()
	go handleUpdates()
	go handleHandshakes(conn)
	listen(conn)
	// for {
	// 	// print("1")
	// 	buf := make([]byte, 1024)
	// 	// print("2")
	// 	_, addr, err := conn.ReadFrom(buf)
	// 	if _, existingClient := clients[addr.String()]; !existingClient {
	// 		print("Launching publish")
	// 		go publish(conn, addr)
	// 		clients[addr.String()] = true
	// 	}
	// 	// print("3")
	// 	if err != nil {
	// 		continue
	// 	}
	// 	// print("4")
	// 	// go serve(conn, addr, buf[:n])
	// }
}

//TODO need to implement a handshake and block clients that don't abide
func listen(conn net.PacketConn) {
	passphrase := []byte("passphrase")
	for {
		t := &state.Tank{}
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buf)
		// print("\n")
		if err != nil {
			continue
		}
		print(fmt.Sprintf("%s", hex.Dump(buf[:n])))
		if bytes.Compare(buf[:n], passphrase) == 0 {
			print("got passphrase. sending to handshake queue\n")
			handshakeQueue <- &addr
		} else {
			state.TankFromBytes(buf[:n], t)
			// fmt.Printf("Enqueue %+v\n", t)
			listenQueue <- t
		}
		// state.ClientStateFromBytes(buf[:n], cs)
		// listenQueue.PushBack(cs)
	}
}

func handleHandshakes(conn net.PacketConn) {
	print("handleHandshakes")
	buf := make([]byte, 4)
	for {
		addr := <-handshakeQueue
		print("processing handshake...\n")
		binary.LittleEndian.PutUint32(buf, uint32(nextClientId))
		print(fmt.Sprintf("sending: %d\n", nextClientId))
		conn.WriteTo(buf, *addr)

		nextClientId++
	}
}

// func handleUpdates(c chan *state.Tank) {
func handleUpdates() {
	// print("handleUpdates")
	for {
		t := <-listenQueue
		//check if t contains the secret phrase
		//	if yes, do handshaek
		//	if no, is there a client id that matches our know list
		print(fmt.Sprintf("handle update: %v\n", t))
	}
}

func dropDisconnectedClients() {
	done := make(chan bool, 1)
	ticker := time.NewTicker(16667 * time.Microsecond) //~60 fps
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// now := time.Now()
			for _, client := range clients {
				if time.Since(client.LastHeardFrom()).Milliseconds() > clientTimeoutMillis {
					print(fmt.Sprintf("Client %d timed out\n", client))
				}
			}
		}
	}
}

// func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
// 	//returns the same data back over the wire but with the 3rd byet modified and the first byte set to 65 (A)
// 	// print("serving\n")
// 	//print(addr.Network())
// 	//print(addr.String())
// 	num := binary.BigEndian.Uint64(buf)
// 	print(num)
// 	// 0 - 1: ID
// 	// 2: QR(1): Opcode(4)
// 	buf[2] |= 0x80 // Set QR bit
// 	buf[0] = 65
// 	print(fmt.Sprintf("serve Bytes: %s\n", buf))
// 	print(fmt.Sprintf("serve Addr: %s\n", addr))
// 	pc.WriteTo(buf, addr)
// 	for i := 65; i< 75; i++ {
// 		binary.BigEndian.PutUint64(buf, uint64(i))
// 		print(fmt.Sprintf("serve Bytes: %s i %d \n", buf, i))
// 		pc.WriteTo(buf, addr)
// 		time.Sleep(time.Second)
// 	}

// }

//connect to this with
//nc -u localhost 1984
// func publish(pc net.PacketConn, addr net.Addr){

// 	print("Inside publish\n")
// 	bytes := make([]byte, 8)
// 	print(fmt.Sprintf("publish Bytes: %s\n", bytes))
// 	print(fmt.Sprintf("publish Addr: %s\n", addr))
// 	pc.WriteTo(bytes, addr)
// 	for i := 65; i< 75; i++ {
// 		binary.BigEndian.PutUint64(bytes, uint64(i))
// 		print(fmt.Sprintf("publish Bytes: %s i %d \n", bytes, i))
// 		pc.WriteTo(bytes, addr)
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }
