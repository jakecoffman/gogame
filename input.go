package gogame

import "github.com/hajimehoshi/ebiten"

type Input struct {
	keyState map[ebiten.Key]int
}

func NewInput() *Input {
	return &Input{
		keyState: map[ebiten.Key]int{},
	}
}

type Dir int

var keyMap = []ebiten.Key{
	ebiten.KeyUp,
	ebiten.KeyRight,
	ebiten.KeyDown,
	ebiten.KeyLeft,

	ebiten.KeyEscape,
}

func (i *Input) Update() {
	for _, k := range keyMap {
		if ebiten.IsKeyPressed(k) {
			i.keyState[k] = 1
		} else {
			i.keyState[k] = 0
		}
	}
}
