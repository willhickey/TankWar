package main

import (
	"encoding/hex"
	"fmt"
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