package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/willhickey/TankWar/images"
	"github.com/willhickey/TankWar/network"

)

const (
	screenWidth  = 1600
	screenHeight = 960
	tankWidth = 32
	tankHeight = 32
)

// var (
// 	ebitenImage *ebiten.Image
// )

type Tank struct {
	x        float64
	y        float64
	acc      int
	maxSpeed int
	speed    int
	angle    int
}

func (t *Tank) Update() {
	t.UpdateSpeed()
	t.UpdateAngle()
	t.UpdatePosition()
}

func (t *Tank) UpdatePosition() {
	newX := t.x + float64(t.speed) * math.Cos(float64(t.angle)*math.Pi/float64(180.0))
	newY := t.y + float64(t.speed) * math.Sin(float64(t.angle)*math.Pi/float64(180.0))
	if newX >= 0 && newX <= screenWidth-tankWidth && newY >=0 && newY <= screenHeight - tankHeight {
		//todo if the tank speed is -3, and t.x is 2, then it won't proceed to the wall, it stops short
		t.x = newX
		t.y = newY
	}
}

func (t *Tank) UpdateSpeed() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		print("Key up")
		t.speed = t.speed + t.acc
		if t.speed > t.maxSpeed {
			t.speed = t.maxSpeed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
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

func (t *Tank) UpdateAngle() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		t.angle -= 2
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		t.angle += 2
	}
}

var (
	tankImg *ebiten.Image
	op      = &ebiten.DrawImageOptions{}
	myTank  *Tank
)

func init() {
	myTank = &Tank{
		x:        20,
		y:        20,
		acc:      1,
		maxSpeed: 8,
		speed:    0,
		angle:    0,
	}

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

// func leftTouched() bool {
// 	for _, id := range ebiten.TouchIDs() {
// 		x, _ := ebiten.TouchPosition(id)
// 		if x < screenWidth/2 {
// 			return true
// 		}
// 	}
// 	return false
// }

// func rightTouched() bool {
// 	for _, id := range ebiten.TouchIDs() {
// 		x, _ := ebiten.TouchPosition(id)
// 		if x >= screenWidth/2 {
// 			return true
// 		}
// 	}
// 	return false
// }

func update(screen *ebiten.Image) error {
	// if ebiten.IsKeyPressed(ebiten.KeyUp) {
	// 	myTank.y = myTank.y - 2
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyDown) {
	// 	myTank.y = myTank.y + 2
	// }
	myTank.Update()
	//print(fmt.Sprintf("%d,%d\n", myTank.x, myTank.y))
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.RGBA{0x80, 0x80, 0xc0, 0xff})
	ebitenutil.DrawRect(screen, 0, 0, 5, screenHeight, color.RGBA{255, 0, 0, 150})
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, 5, color.RGBA{255, 0, 0, 150})
	ebitenutil.DrawRect(screen, screenWidth-5, 0, 5, screenHeight, color.RGBA{255, 0, 0, 150})
	ebitenutil.DrawRect(screen, 0, screenHeight-5, screenWidth, 5, color.RGBA{255, 0, 0, 150})
	op.GeoM.Reset()
	w, h := tankImg.Size()

	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(2 * math.Pi * float64(myTank.angle) / 360)
	op.GeoM.Translate(float64(w)/2, float64(h)/2)

	op.GeoM.Translate(float64(myTank.x), float64(myTank.y))
	screen.DrawImage(tankImg, op)
	return nil
}

func main() {
	ebiten.SetFullscreen(true)
	// w, h := ebiten.ScreenSizeInFullscreen()
	// s := ebiten.DeviceScaleFactor()
	print(fmt.Sprintf("%d,%d\n", myTank.x, myTank.y))
	print(fmt.Sprintf("%d\n", network.RxQueue.Len()))
	network.Foo()
	print(fmt.Sprintf("%d\n", network.RxQueue.Len()))

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "TankWar"); err != nil {
		log.Fatal(err)
	}
}
