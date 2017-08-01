package gogame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"fmt"
	"net"
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
	ID      int8
	IsLocal bool
	addr    *net.UDPAddr `gob:"-"`

	Shape *chipmunk.Shape

	Image *ebiten.Image
}

func NewPlayer(isLocal bool, addr *net.UDPAddr) *Player {
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
		IsLocal: isLocal,
		addr:    addr,
		Image:   square,
		Shape:   box,
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

		var velocity float32 = 0.0
		if input.keyState[ebiten.KeyW] == 1 {
			velocity = maxSpeed * -1
		}
		if input.keyState[ebiten.KeyS] == 1 {
			velocity = maxSpeed
		}

		move := &Move{angularVelocity, velocity}
		Send(move.Serialize(), serverAddr)
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

	ebitenutil.DebugPrint(screen, fmt.Sprint("\nPlayer ", p.Shape.Body.Position()))
}
