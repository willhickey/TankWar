package main

import (
	// "log"
	"fmt"
	// "os"
	"net"
	"github.com/pkg/term"
	"encoding/binary"
)


func main() {
	//CONNECT := os.Args[1]
	s, err := net.ResolveUDPAddr("udp4", "localhost:1984") //CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 1024)
	n, addr, err := c.ReadFrom(buf)


	defer c.Close()
	for {
		ascii, _, err := getChar()
		if err != nil {
			return
		}
		print(ascii)
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, uint64(ascii))
		_, err = c.Write(bytes)
		
		if ascii == 17 || ascii == 3 {
			print("\n")
			break
		}

	}
}



func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}