package state

import (
	"fmt"
	"encoding/hex"
	"net"
	"testing"
	"time"

)

func TestToBytes(t *testing.T) {
	buf := make([]byte, 13)

	myTank := &Tank{
		x:        20,
		y:        22,
		acc:      1,
		maxSpeed: 8,
		speed:    -1,
		angle:    256,
	}
	myTank2 := &Tank{}

	TankToBytes(myTank, buf)
	TankFromBytes(buf, myTank2)

	fmt.Printf("%+v\n", myTank)
	print(fmt.Sprintf("%s", hex.Dump(buf)))
	fmt.Printf("%+v\n", myTank2)
}

func TestClientStateToBytes(t *testing.T) {
	//make a ClientState
	myTank := NewTank(20, 20, 1, 8, 0, 0)
	q := net.ParseIP("192.168.0.255")
	addr := &net.IPAddr{q,"udp"}
	myClientState := NewClientState(13, addr, time.Now(), myTank)
	print(fmt.Sprintf("%d\n%v\n%s\n%+v\n", myClientState.id, myClientState.address, myClientState.lastHeardFrom.Format(time.RFC3339), myClientState.tank))
	myClientStateBytes := make([]byte, 100)
	ClientStateToBytes(myClientState, myClientStateBytes)
	//To Bytes it
}