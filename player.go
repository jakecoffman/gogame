package gogame

import (
	"fmt"
	"image/color"
	"net"

	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

const (
	startX       = 128.0
	startY       = 128.0
	playerWidth  = 32.0
	playerHeight = 32.0

	maxSpeed  float32 = 100.0
	maxTorque float32 = 0.05
)

type Player struct {
	ID   int8
	Addr *net.UDPAddr

	Shape *chipmunk.Shape

	Image *ebiten.Image
}

func (p *Player) IsLocal() bool {
	return !IsServer && p.ID == Me
}

func NewPlayer() *Player {
	square, _ := ebiten.NewImage(playerWidth, playerHeight, ebiten.FilterNearest)

	// chipmunk origin is the bottom left corner
	box := chipmunk.NewBox(vect.Vect{0, 0}, playerWidth, playerHeight)
	box.SetElasticity(0.5)
	box.SetFriction(5.0)

	body := chipmunk.NewBody(1.0, box.Moment(1.0))
	body.SetPosition(vect.Vect{startX, startY})
	//body.SetAngle(vect.Float(rand.Float32() * 2 * math.Pi))
	body.SetMass(1.0)

	body.AddShape(box)
	space.AddBody(body)

	return &Player{
		Image: square,
		Shape: box,
	}
}

func (p *Player) Location() *Location {
	return &Location{
		ID:              p.ID,
		X:               float32(p.Shape.Body.Position().X),
		Y:               float32(p.Shape.Body.Position().Y),
		Angle:           float32(p.Shape.Body.Angle()),
		AngularVelocity: float32(p.Shape.Body.AngularVelocity()),
		Vx:              float32(p.Shape.Body.Velocity().X),
		Vy:              float32(p.Shape.Body.Velocity().Y),
	}
}

var lastTurn, lastThrottle float32

func (p *Player) Update() {
	if p.IsLocal() {
		var turn float32
		if Input.keyState[ebiten.KeyA] == 1 {
			turn = maxTorque * -1
		}
		if Input.keyState[ebiten.KeyD] == 1 {
			turn = maxTorque
		}

		var throttle float32
		if Input.keyState[ebiten.KeyW] == 1 {
			throttle = maxSpeed * -1
		}
		if Input.keyState[ebiten.KeyS] == 1 {
			throttle = maxSpeed
		}

		if lastTurn != turn || lastThrottle != throttle {
			lastTurn = turn
			lastThrottle = throttle
			move := &Move{Turn: turn, Throttle: throttle}
			bin, err := move.MarshalBinary()
			if err != nil {
				log.Println(err)
				return
			}
			Send(bin, ServerAddr)
		}
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

	if p.IsLocal() {
		ebitenutil.DebugPrint(screen, fmt.Sprint("\nPlayer ", p.Shape.Body.Position(), "\n", ))
	}
}
