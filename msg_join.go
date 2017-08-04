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
		b, err := (&Join{ID: player.ID, You: true}).MarshalBinary()
		if err != nil {
			log.Println(err)
			return err
		}
		Send(b, addr)
		loc, err := player.Location().MarshalBinary()
		if err != nil {
			log.Println(err)
			return err
		}
		// tell this player where they are
		Send(loc, addr)
		joinBytes, err := Join{player.ID, false}.MarshalBinary()
		if err != nil {
			log.Println(err)
			return err
		}
		for _, p := range Players {
			// tell all players about this player
			Send(joinBytes, p.Addr)
			Send(loc, p.Addr)
			// tell this player where all the existing players are
			b, err = (&Join{p.ID, false}).MarshalBinary()
			if err != nil {
				log.Println(err)
				continue
			}
			Send(b, player.Addr)
			b, err = p.Location().MarshalBinary()
			if err != nil {
				log.Println(err)
				continue
			}
			Send(b, player.Addr)
		}
	} else {
		log.Println("Player joined")
		player.ID = j.ID
		if j.You {
			log.Println("Oh, it's me!")
			Me = player.ID
			// now that I am joined I will start pinging the server
			go PingRegularly()
		}
	}
	Players[player.ID] = player

	return nil
}

func (j Join) MarshalBinary() ([]byte, error) {
	if j.You {
		return []byte{JOIN, byte(j.ID), 1}, nil
	} else {
		return []byte{JOIN, byte(j.ID), 0}, nil
	}
}

func (j *Join) UnmarshalBinary(b []byte) error {
	j.ID = int8(b[1])
	if b[2] == 0 {
		j.You = false
	} else {
		j.You = true
	}
	return nil
}
