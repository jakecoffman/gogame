package gogame

import (
	"fmt"
	"image/color"
	"log"

	"errors"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/vova616/chipmunk"
)

var space *chipmunk.Space

var input *Input
var players []*Player

const size = 400

func update(screen *ebiten.Image) error {
	screen.Fill(color.NRGBA{0x00, 0x00, 0x00, 0xff})

	input.Update()

	if input.keyState[ebiten.KeyEscape] == 1 {
		return errors.New("Player quit")
	}

	for _, player := range players {
		player.Update()
	}

	space.Step(1.0/60.0)

	if ebiten.IsRunningSlowly() {
		return nil
	}

	DrawLevel(screen)

	for _, player := range players {
		player.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%f", ebiten.CurrentFPS()))
	return nil
}

func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Game starting")
	defer func() { log.Println("Game ended") }()

	input = NewInput()

	LevelInit()

	players = []*Player{NewPlayer(true)}

	if err := ebiten.Run(update, size, size, 1, "Hello, world!"); err != nil {
		log.Fatal(err)
	}
}
