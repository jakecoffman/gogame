package gogame

import (
	"net"
	"log"
	"os"
	"encoding/gob"
	"bytes"
	"github.com/vova616/chipmunk"
	"math"
)

var serverAddr *net.UDPAddr
var udpConn *net.UDPConn
var isServer bool

type Incoming struct {
	data []byte
	addr    *net.UDPAddr
}

var incomings chan Incoming

func init() {
	serverAddr = &net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	incomings = make(chan Incoming)
}

func NetInit() {
	var err error
	if len(os.Args) > 1 && os.Args[1] == "server" {
		isServer = true
	}

	if isServer {
		udpConn, err = net.ListenUDP("udp", serverAddr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		udpConn, err = net.DialUDP("udp", nil, serverAddr)
		if err != nil {
			log.Fatal(err)
		}
	}

	udpConn.SetReadBuffer(1048576)

	go Recv()
}

func NetClose() error {
	return udpConn.Close()
}

func Recv() {
	for {
		data := make([]byte, 2048)
		_, addr, err := udpConn.ReadFromUDP(data)
		if err != nil {
			panic(err)
		}
		incomings <- Incoming{data, addr}
		log.Println("Message queued")
	}
}

func Process() {
	var err error
	for {
		select {
		case incoming := <-incomings:
			log.Println("Processing message")

			switch incoming.data[0] {
			case 0:
				if err = (Join{}).Handle(incoming.addr); err != nil {
					log.Fatal(err)
				}
			}
		default:
			// no data to process this frame
			return
		}
	}
}

func Send(data []byte, addr *net.UDPAddr) {
	if isServer {
		_, err := udpConn.WriteToUDP(data, addr)
		if err != nil {
			panic(err)
		}
	} else {
		_, err := udpConn.Write(data)
		if err != nil {
			panic(err)
		}
	}
}

var curId int8

type Join struct {
	ID int8
	You bool
}

func (j Join) Handle(addr *net.UDPAddr) error {
	log.Println("Player has joined")
	player := NewPlayer(false, addr)

	if isServer {
		player.ID = curId
		curId++
		Send(Join{ID: player.ID, You: true}.Marshal(), addr)
		for _, p := range players {
			join := &Join{player.ID, false}
			Send(join.Marshal(), p.addr)
		}
		players[player.ID] = player
	}

	return nil
}

func (j Join) Marshal() []byte {
	if j.You {
		return []byte{0, byte(j.ID), 1}
	} else {
		return []byte{0, byte(j.ID), 0}
	}
}

func (j Join) Unmarshal(b []byte) {
	j.ID = int8(b[1])
	if b[2] == 0 {
		j.You = false
	} else {
		j.You = true
	}
}

//type You struct {
//	ID int8
//}
//
//func (y You) Handle(addr *net.UDPAddr) error {
//	log.Println("Handling YOU")
//	player := NewPlayer(true, nil)
//	player.ID = y.ID
//	return nil
//}
//
//func (y You) Serialize() []byte {
//	var bytez []byte
//	buf := bytes.NewBuffer(bytez)
//	err := gob.NewEncoder(buf).Encode(y)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return bytez
//}

type Move struct {
	AngularVelocity, Velocity float32
}

func (m Move) Handle(addr *net.UDPAddr) error {
	log.Println("Player has moved")

	// handle movement on the server and send update to players
	var player *Player
	for _, p := range players {
		if p.addr.String() == addr.String() {
			player = p
		}
	}

	player.Shape.Body.SetAngularVelocity(m.AngularVelocity)
	vx2 := math.Cos(float64(player.Shape.Body.Angle() * chipmunk.DegreeConst))
	vy2 := math.Sin(float64(player.Shape.Body.Angle() * chipmunk.DegreeConst))
	svx2, svy2 := m.Velocity*float32(vx2), m.Velocity*float32(vy2)
	player.Shape.Body.SetVelocity(svx2, svy2)

	for _, p := range players {
		Send(m.Serialize(), p.addr)
	}

	return nil
}

func (m Move) Serialize() []byte {
	var bytez []byte
	buf := bytes.NewBuffer(bytez)
	err := gob.NewEncoder(buf).Encode(m)
	if err != nil {
		log.Fatal(err)
	}
	return bytez
}
