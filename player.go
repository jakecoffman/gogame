package gogame

import (
	"fmt"
	"image/color"
	"log"
	"net"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/jakecoffman/physics"
)

const (
	startX       = 20.0
	startY       = 20.0
	playerWidth  = 32.0
	playerHeight = 32.0
)

type Player struct {
	ID   int8
	Addr *net.UDPAddr

	Shape *physics.Shape

	Image *ebiten.Image
}

func (p *Player) IsLocal() bool {
	return !IsServer && p.ID == Me
}

func NewPlayer() *Player {
	square, _ := ebiten.NewImage(playerWidth, playerHeight, ebiten.FilterNearest)

	radius := (&physics.Vector{playerWidth, playerHeight}).Length()
	body := space.AddBody(physics.NewBody(1, physics.MomentForBox(1, playerWidth, playerHeight)))
	body.SetPosition(&physics.Vector{startX, startY})
	shape := space.AddShape(physics.NewBox(body, playerWidth, playerHeight, radius))
	shape.E = 0
	shape.U = 5

	return &Player{
		Image: square,
		Shape: shape,
	}
}

func (p *Player) Location() *Location {
	return &Location{
		ID:              p.ID,
		X:               p.Shape.Body().Position().X,
		Y:               p.Shape.Body().Position().Y,
		Angle:           p.Shape.Body().Angle(),
		AngularVelocity: p.Shape.Body().AngularVelocity(),
		Vx:              p.Shape.Body().Velocity().X,
		Vy:              p.Shape.Body().Velocity().Y,
	}
}

func (p *Player) Update() {
	if p.IsLocal() {
		var turn float64
		if Input.keyState[ebiten.KeyA] == 1 {
			turn = -1
		}
		if Input.keyState[ebiten.KeyD] == 1 {
			turn = 1
		}

		var throttle float64
		if Input.keyState[ebiten.KeyW] == 1 {
			throttle = -1
		}
		if Input.keyState[ebiten.KeyS] == 1 {
			throttle = 1
		}

		bin, err := (&Move{Turn: turn, Throttle: throttle}).MarshalBinary()
		if err != nil {
			log.Println(err)
			return
		}
		Send(bin, ServerAddr)
	}
}

var opts *ebiten.DrawImageOptions

func (p *Player) Draw(screen *ebiten.Image) {
	p.Image.Fill(color.White)
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-playerWidth/2, -playerHeight/2)
	opts.GeoM.Rotate(p.Shape.Body().Angle() * physics.DegreeConst)
	pos := p.Shape.BB().Center()
	opts.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(p.Image, opts)

	if p.IsLocal() {
		str := fmt.Sprintf("\nPlayer %v\nRot %v\nPing %v",
			p.Shape.Body().Position(), p.Shape.Body().Angle(), LastPing.Get())
		ebitenutil.DebugPrint(screen, str)
	}
}
