package gogame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
)

const (
	startX = 128.0
	startY = 128.0
	playerWidth = 32.0
	playerHeight = 32.0

	maxSpeed  float32 = 100.0
	maxTorque float32 = 0.05
)

type Player struct {
	IsLocal bool

	Shape *chipmunk.Shape

	Image *ebiten.Image
}

func NewPlayer(isLocal bool) *Player {
	square, _ := ebiten.NewImage(playerWidth, playerHeight, ebiten.FilterNearest)

	// chipmunk origin is the bottom left corner
	box := chipmunk.NewBox(vect.Vect{startX, startY}, playerWidth, playerHeight)
	box.SetElasticity(0.5)
	box.SetFriction(5.0)

	body := chipmunk.NewBody(1.0, box.Moment(1.0))
	body.SetPosition(vect.Vect{startX, startY})
	//body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))
	body.SetMass(1.0)

	body.AddShape(box)
	space.AddBody(body)

	return &Player{
		IsLocal: isLocal,
		Image:   square,
		Shape: box,
	}
}

func (p *Player) Update() {
	if p.IsLocal {
		var angularVelocity float32 = 0.0
		if input.keyState[ebiten.KeyA] == 1 {
			angularVelocity = maxTorque * -1
		}
		if input.keyState[ebiten.KeyD] == 1 {
			angularVelocity = maxTorque
		}

		p.Shape.Body.SetAngularVelocity(angularVelocity)

		var velocity float32 = 0.0
		if input.keyState[ebiten.KeyW] == 1 {
			velocity = maxSpeed * -1
		}
		if input.keyState[ebiten.KeyS] == 1 {
			velocity = maxSpeed
		}

		vx2 := math.Cos(float64(p.Shape.Body.Angle() * chipmunk.DegreeConst))
		vy2 := math.Sin(float64(p.Shape.Body.Angle() * chipmunk.DegreeConst))
		svx2, svy2 := velocity * float32(vx2), velocity * float32(vy2)
		p.Shape.Body.SetVelocity(svx2, svy2)
	}
}

var opts *ebiten.DrawImageOptions

func (p *Player) Draw(screen *ebiten.Image) {
	p.Image.Fill(color.White)
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-playerWidth/2, -playerHeight/2)
	opts.GeoM.Rotate(float64(p.Shape.Body.Angle() * chipmunk.DegreeConst))
	opts.GeoM.Translate(float64(p.Shape.Body.Position().X), float64(p.Shape.Body.Position().Y))
	screen.DrawImage(p.Image, opts)

	ebitenutil.DebugPrint(screen, fmt.Sprint("\nPlayer ", p.Shape.Body.Position(), p.Shape.BB.Upper, ",", p.Shape.BB.Lower))
}
