package main

import (
	"log"
	"runtime/pprof"

	"os"

	"os/signal"

	"github.com/hajimehoshi/ebiten"
	"github.com/jakecoffman/gogame"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
			os.Exit(1)
		}
	}()

	log.Println("Game starting")
	defer func() { log.Println("Game ended") }()

	gogame.NetInit()
	defer func() { log.Println(gogame.NetClose()) }()

	gogame.Input = gogame.NewInput()
	gogame.LevelInit()

	join := &gogame.Join{}
	log.Println("Sending JOIN command")

	b, err := join.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	gogame.Send(b, gogame.ServerAddr)

	ebiten.SetRunnableInBackground(true)

	if err := ebiten.Run(gogame.Update, gogame.Size, gogame.Size, 1, "Client"); err != nil {
		log.Fatal(err)
	}
}
