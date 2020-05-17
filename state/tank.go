package state

import (
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	ScreenWidth  = 1600
	ScreenHeight = 960
	TankWidth = 32
	TankHeight = 32
	TankByteLength = 13
)

type Tank struct {
	x        float32
	y        float32
	acc      int8
	maxSpeed int8
	speed    int8
	angle    int16
	//TODO add client frame number for server packet ordering
}

func NewTank(x float32, y float32, acc int8, maxSpeed int8, speed int8, angle int16) *Tank {
	t := new(Tank)
	t.x = x
	t.y = y
	t.acc = acc
	t.maxSpeed = maxSpeed
	t.speed = speed
	t.angle = angle
	return t
}

func (t Tank) X() float32 {
	return t.x
}
func (t Tank) Y() float32 {
	return t.y
}
func (t Tank) Angle() int16 {
	return t.angle
}

func (t *Tank) Update(keyState map[ebiten.Key]bool) {
	t.UpdateSpeed(keyState)
	t.UpdateAngle(keyState)
	t.UpdatePosition()
}

func (t *Tank) UpdatePosition() { //screenWidth int, screenHeight int) {
	newX := t.x + float32(t.speed) * float32(math.Cos(float64(t.angle)*math.Pi/float64(180.0)))
	newY := t.y + float32(t.speed) * float32(math.Sin(float64(t.angle)*math.Pi/float64(180.0)))
	if newX >= 0 && newX <= ScreenWidth-TankWidth && newY >=0 && newY <= ScreenHeight - TankHeight {
		//todo if the tank speed is -3, and t.x is 2, then it won't proceed to the wall, it stops short
		t.x = newX
		t.y = newY
	}
}

func (t *Tank) UpdateSpeed(keyState map[ebiten.Key]bool) {
	if keyState[ebiten.KeyUp] {
		t.speed = t.speed + t.acc
		if t.speed > t.maxSpeed {
			t.speed = t.maxSpeed
		}
	} else if keyState[ebiten.KeyDown] {
		t.speed = t.speed - t.acc
		if t.speed < -1*t.maxSpeed {
			t.speed = -1 * t.maxSpeed
		}
	} else { //coasting
		if t.speed > 0 {
			t.speed = t.speed - t.acc
			if t.speed < 0 {
				t.speed = 0
			}
		} else if t.speed < 0 {
			t.speed = t.speed + t.acc
			if t.speed > 0 {
				t.speed = 0
			}
		}
	}
}

func (t *Tank) UpdateAngle(keyState map[ebiten.Key]bool) {
	if keyState[ebiten.KeyLeft] {
		t.angle -= 2
	} else if keyState[ebiten.KeyRight] {
		t.angle += 2
	}
}


func TanksToBytes(tanks []*Tank, buf []byte) {

}

func TankToBytes(t *Tank, buf []byte) {
	float32ToBytes(t.x, buf[0:])
	float32ToBytes(t.y, buf[4:])
	buf[8] = byte(t.acc)
	buf[9] = byte(t.maxSpeed)
	buf[10] = byte(t.speed)
	binary.LittleEndian.PutUint16(buf[11:], uint16(t.angle))
}

func TankFromBytes(buf []byte, t *Tank) {
	t.x = float32FromBytes(buf[0:])
	t.y = float32FromBytes(buf[4:])
	t.acc = int8(buf[8])
	t.maxSpeed = int8(buf[9])
	t.speed = int8(buf[10])
	t.angle = int16(binary.LittleEndian.Uint16(buf[11:]))
}

func float32FromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func float32ToBytes(f float32, buf []byte){
	bits := math.Float32bits(f)
	//bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, bits)
}
