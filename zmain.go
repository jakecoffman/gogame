package gogame

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/vova616/chipmunk"
)

var space *chipmunk.Space

var Input *input
var Players map[int8]*Player

// Server only lookup of addr to ID
var Lookup map[string]int8
var Me int8

const Size = 400

func init() {
	Players = map[int8]*Player{}
	Lookup = map[string]int8{}
}

func Update(screen *ebiten.Image) error {
	Process()

	if !IsServer {
		screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})

		Input.Update()

		if Input.keyState[ebiten.KeyEscape] == 1 {
			return errors.New("Player quit")
		}
		if Input.keyState[ebiten.KeyF10] == 1 {
			panic("User invoked crash")
		}
	}

	for _, player := range Players {
		player.Update()
	}

	space.Step(1.0 / 60.0)

	if !IsServer {
		if ebiten.IsRunningSlowly() {
			return nil
		}

		DrawLevel(screen)

		for _, player := range Players {
			player.Draw(screen)
		}

		ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.CurrentFPS()))
	}
	return nil
}
