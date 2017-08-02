package gogame

import (
	"fmt"
	"image/color"
	"net"

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
	addr *net.UDPAddr

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

var lastAngularVelocity, lastVelocity float32

func (p *Player) Update() {
	if p.IsLocal() {
		var angularVelocity float32
		if Input.keyState[ebiten.KeyA] == 1 {
			angularVelocity = maxTorque * -1
		}
		if Input.keyState[ebiten.KeyD] == 1 {
			angularVelocity = maxTorque
		}

		var velocity float32
		if Input.keyState[ebiten.KeyW] == 1 {
			velocity = maxSpeed * -1
		}
		if Input.keyState[ebiten.KeyS] == 1 {
			velocity = maxSpeed
		}

		if lastAngularVelocity != angularVelocity || lastVelocity != velocity {
			lastAngularVelocity = angularVelocity
			lastVelocity = velocity
			move := &Move{AngularVelocity: angularVelocity, Velocity: velocity}
			Send(move.Marshal(), ServerAddr)
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

	ebitenutil.DebugPrint(screen, fmt.Sprint("\nPlayer ", p.Shape.Body.Position()))
}
