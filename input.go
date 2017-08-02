package gogame

import "github.com/hajimehoshi/ebiten"

type input struct {
	keyState map[ebiten.Key]int
}

func NewInput() *input {
	return &input{
		keyState: map[ebiten.Key]int{},
	}
}

type Dir int

var keyMap = []ebiten.Key{
	ebiten.KeyUp,
	ebiten.KeyRight,
	ebiten.KeyDown,
	ebiten.KeyLeft,

	ebiten.KeyW,
	ebiten.KeyA,
	ebiten.KeyS,
	ebiten.KeyD,

	ebiten.KeyEscape,
	ebiten.KeyF10,
}

func (i *input) Update() {
	for _, k := range keyMap {
		if ebiten.IsKeyPressed(k) {
			i.keyState[k] = 1
		} else {
			i.keyState[k] = 0
		}
	}
}
