package gogame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type Player struct {
	IsLocal bool
	Pos     Vec

	Image *ebiten.Image
}

func NewPlayer(isLocal bool) *Player {
	square, _ := ebiten.NewImage(32, 32, ebiten.FilterNearest)
	return &Player{
		IsLocal: isLocal,
		Image:   square,
		Pos:     Vec{64, 64},
	}
}

func (p *Player) Update() {
	if p.IsLocal {
		p.Pos.x += float64(input.keyState[ebiten.KeyRight] * movementSpeed)
		p.Pos.x -= float64(input.keyState[ebiten.KeyLeft] * movementSpeed)
		p.Pos.y -= float64(input.keyState[ebiten.KeyUp] * movementSpeed)
		p.Pos.y += float64(input.keyState[ebiten.KeyDown] * movementSpeed)
	}
}

var opts *ebiten.DrawImageOptions

func (p *Player) Draw(screen *ebiten.Image) {
	p.Image.Fill(color.White)
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.Pos.x, p.Pos.y)
	screen.DrawImage(p.Image, opts)
}
