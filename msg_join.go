package gogame

import (
	"net"
	"log"
)

var curId int8 = 1

type Join struct {
	ID  int8
	You bool
}

func (j *Join) Handle(addr *net.UDPAddr) error {
	player := NewPlayer()

	if IsServer {
		// player initialization (TODO set spawn point)
		player.Addr = addr
		player.ID = curId
		curId++
		Lookup[addr.String()] = player.ID
		// tell this player their ID
		Send((&Join{ID: player.ID, You: true}).Marshal(), addr)
		loc := player.Location().Marshal()
		// tell this player where they are
		Send(loc, addr)
		joinMsg := Join{player.ID, false}.Marshal()
		for _, p := range Players {
			// tell all players about this player
			Send(joinMsg, p.Addr)
			Send(loc, p.Addr)
			// tell this player where all the existing players are
			Send(Join{p.ID, false}.Marshal(), player.Addr)
			Send(p.Location().Marshal(), player.Addr)
		}
	} else {
		log.Println("Player joined")
		player.ID = j.ID
		if j.You {
			log.Println("Oh, it's me!")
			Me = player.ID
		}
	}
	Players[player.ID] = player

	return nil
}

func (j Join) Marshal() []byte {
	if j.You {
		return []byte{JOIN, byte(j.ID), 1}
	} else {
		return []byte{JOIN, byte(j.ID), 0}
	}
}

func (j *Join) Unmarshal(b []byte) {
	j.ID = int8(b[1])
	if b[2] == 0 {
		j.You = false
	} else {
		j.You = true
	}
}
