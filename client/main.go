package main

import (
	"bytes"
	"fmt"
	"encoding/hex"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"net"

	"github.com/willhickey/TankWar/state"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/willhickey/TankWar/images"
)

// const (
// 	screenWidth  = 1600
// 	screenHeight = 960
// 	tankWidth = 32
// 	tankHeight = 32
// )

// var (
// 	ebitenImage *ebiten.Image
// )

// type Tank struct {
// 	x        float64
// 	y        float64
// 	acc      int
// 	maxSpeed int
// 	speed    int
// 	angle    int
// }

var (
	tankImg *ebiten.Image
	op      = &ebiten.DrawImageOptions{}
	myTank  *state.Tank
	keyState = make(map[ebiten.Key]bool)
	myTankBytes = make([]byte, state.TankByteLength)
	keys = [...]ebiten.Key{ebiten.KeyUp,
		ebiten.KeyDown,
		ebiten.KeyLeft,
		ebiten.KeyRight}
	done = make(chan bool)
	clientId int32 = -1
)

func init() {
	myTank = state.NewTank(20, 20, 1, 8, 0, 0)
	// &state.Tank{
	// 	x:        20,
	// 	y:        20,
	// 	acc:      1,
	// 	maxSpeed: 8,
	// 	speed:    0,
	// 	angle:    0,
	// }
	img, _, err := image.Decode(bytes.NewReader(images.Tank_png))
	if err != nil {
		log.Fatal(err)
	}
	tankImg, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	// w, h := tankImg.Size()
	// ebitenImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)

	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.5)
	op.GeoM.Rotate(2 * math.Pi * float64(100) / 360)
	//ebitenImage.DrawImage(tankImg, op)

}

func getKeyState() {

	for _, key := range keys {
		keyState[key] = ebiten.IsKeyPressed(key)
	}
}

func update(screen *ebiten.Image) error {

	getKeyState()
	myTank.Update(keyState)
	state.TankToBytes(myTank, myTankBytes)
	print(fmt.Sprintf("%s", hex.Dump(myTankBytes)))

	// print(fmt.Sprintf("%d,%d\n", myTank.x, myTank.y))
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.RGBA{0x80, 0x80, 0xc0, 0xff})
	ebitenutil.DrawRect(screen, 0, 0, 5, state.ScreenHeight, color.RGBA{255, 0, 0, 150})
	ebitenutil.DrawRect(screen, 0, 0, state.ScreenWidth, 5, color.RGBA{255, 0, 0, 150})
	ebitenutil.DrawRect(screen, state.ScreenWidth-5, 0, 5, state.ScreenHeight, color.RGBA{255, 0, 0, 150})
	ebitenutil.DrawRect(screen, 0, state.ScreenHeight-5, state.ScreenWidth, 5, color.RGBA{255, 0, 0, 150})
	op.GeoM.Reset()
	w, h := tankImg.Size()

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(2 * math.Pi * float64(myTank.Angle()) / 360)
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(float64(myTank.X()), float64(myTank.Y()))
	screen.DrawImage(tankImg, op)
	return nil
}

func main() {

	

	ebiten.SetFullscreen(true)
	// w, h := ebiten.ScreenSizeInFullscreen()
	// s := ebiten.DeviceScaleFactor()
	
	/*
	connect to server
	spawn listener thread
	spawn publisher thread
	*/


	s, _ := net.ResolveUDPAddr("udp4", "localhost:1984") //CONNECT)
	c, _ := net.DialUDP("udp4", nil, s)
	//TODO handshake
	clientId, _ = doHandshake(c)

	go publish(c)

	if err := ebiten.Run(update, state.ScreenWidth, state.ScreenHeight, 1, "TankWar"); err != nil {
		log.Fatal(err)
	}
	done <- true
}
