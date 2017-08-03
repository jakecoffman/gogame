package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/jakecoffman/gogame"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Game starting")
	defer func() { log.Println("Game ended") }()

	gogame.NetInit()
	defer func() { log.Println(gogame.NetClose()) }()

	gogame.Input = gogame.NewInput()
	gogame.LevelInit()

	join := &gogame.Join{}
	log.Println("Sending JOIN command")

	gogame.Send(join.Marshal(), gogame.ServerAddr)

	ebiten.SetRunnableInBackground(true)

	if err := ebiten.Run(gogame.Update, gogame.Size, gogame.Size, 1, "Client"); err != nil {
		log.Fatal(err)
	}
}
