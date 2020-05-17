package state

import (
	// "encoding/binary"
	// "encoding/hex"
	// "fmt"
	"net"
	"time"

)

const(
	ClientStateByteLength = 100
)


/*
Client sends to server:
	id (from handshake)
	Tank
	clientFrameNumber

Server sends to Client:
	serverFrameNumber
		id
		tank

Should the server send frames for tanks that haven't been heard from in a while? I think no.
Should the server eventually send a "timed out" message for a lost client? Or just let the other clients run the 
tank to the wall?

Client stores:
	it's own ID (from handshake)
	list of all tanks (with ids?)

Server stores:
	server frame number
	client
		id
		address
		last heard from
		tank

*/

type ClientState struct {
	id uint8
	//clientTick int64
	address net.Addr
	lastHeardFrom time.Time
	tank *Tank
}

func (cs *ClientState) LastHeardFrom() (time.Time) {
	return cs.lastHeardFrom
}

func NewClientState(id uint8, address net.Addr, lastHeardFrom time.Time, tank *Tank) *ClientState {
	cs := new(ClientState)
	cs.id = id
	cs.address = address
	cs.lastHeardFrom = lastHeardFrom
	cs.tank = tank
	return cs
}


//id 0
//addr_len 1
//addr 2 to (1 + addr_len)
//time_len 2 + addr_len
//time 3 + addr_len to (2 + addr_len + time_len)
//tank (3 + addr_len + time_len) to (2 + addr_len + time_len + 13

/*for now intentionally leaving out the address since:
	1) the client shouldn't need it
	2) giving it to the client is a privacy breech
	3) it's variable length which makes marshalling harder
*/
func ClientStateToBytes(cs *ClientState, buf []byte) {
	buf[0] = byte(cs.id)
	lastHeardFromBytes, _ := cs.lastHeardFrom.MarshalBinary()
	copy(buf[1:15], lastHeardFromBytes)
	TankToBytes(cs.tank, buf[16:28])
}

func ClientStateFromBytes(buf []byte, cs *ClientState) {
	cs.id = uint8(buf[0])
	cs.lastHeardFrom.UnmarshalBinary(buf[1:15])
	TankFromBytes(buf[16:28], cs.tank)
}

// func FromBytes(buf []byte, cs *ClientState) {
// 	t.x = float32FromBytes(buf[0:])
// 	t.y = float32FromBytes(buf[4:])
// 	t.acc = int8(buf[8])
// 	t.maxSpeed = int8(buf[9])
// 	t.speed = int8(buf[10])
// 	t.angle = int16(binary.LittleEndian.Uint16(buf[11:]))
// }

// type GameState struct {
// 	serverTick int64
// 	tanks []ClientState
// }

// id int //assigned at connection time
// address
// max frame processed
// last heard from
// tank data

/*
buffer length
client ID
tanks []Tank
time / tick / packet?

*/