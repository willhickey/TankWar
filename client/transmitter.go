package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/willhickey/TankWar/state"
)
/*
TODO 
Client sends to server:
	id (from handshake)
	Tank
	clientFrameNumber
*/

func publish(conn *net.UDPConn) {
	clientFrameNumber := 0
	ticker := time.NewTicker(16667 * time.Microsecond)	//~60 fps
	for {
		select {
		case <- done:
			return
		case <-ticker.C:
			state.TankToBytes(myTank, myTankBytes)
			conn.Write(myTankBytes)
			// conn.Write([]byte("\n"))
			print(fmt.Sprintf("Pub %s", hex.Dump(myTankBytes)))
		}
		clientFrameNumber++
	}
}

func doHandshake(conn *net.UDPConn) (int32, error) {
	// send passphrase
	conn.Write([]byte("passphrase"))

	// read a response with a clientid in it
	buf := make([]byte, 4)
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		log.Fatal(err)
	}
	if n != 4 {
		log.Fatal("network error: server returned %d bytes during handshake. expected 4.", n)
	} 
	clientId := int32(binary.LittleEndian.Uint32(buf))
	return clientId, nil 

}