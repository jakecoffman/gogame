package gogame

import (
	"fmt"
	"image/color"
	"log"

	"errors"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var scene int

var input *Input
var players []*Player

const movementSpeed = 1.0

type Vec struct {
	x, y float64
}

func update(screen *ebiten.Image) error {
	// Fill the screen with #FF0000 color
	screen.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})

	input.Update()

	if input.keyState[ebiten.KeyEscape] == 1 {
		return errors.New("Player quit")
	}

	for _, player := range players {
		player.Update()
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

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
	players = []*Player{NewPlayer(true)}
	if err := ebiten.Run(update, 800, 600, 1, "Hello, world!"); err != nil {
		log.Fatal(err)
	}
}
