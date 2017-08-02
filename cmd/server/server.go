package main

import (
	"log"
	"time"

	"github.com/jakecoffman/gogame"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Game starting")
	defer func() { log.Println("Game ended") }()

	gogame.IsServer = true
	gogame.NetInit()
	defer func() { log.Println(gogame.NetClose()) }()

	gogame.LevelInit()

	tick := time.Tick(16 * time.Millisecond)

	for {
		select {
		case <-tick:
			gogame.Update(nil)
		}
	}
}
