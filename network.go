package gogame

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"

	"github.com/vova616/chipmunk"
)

var ServerAddr *net.UDPAddr
var udpConn *net.UDPConn
var IsServer bool

type Incoming struct {
	data []byte
	addr *net.UDPAddr
}

var incomings chan Incoming

func init() {
	ServerAddr = &net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	incomings = make(chan Incoming, 100)
}

func NetInit() {
	var err error

	if IsServer {
		udpConn, err = net.ListenUDP("udp", ServerAddr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		udpConn, err = net.DialUDP("udp", nil, ServerAddr)
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
		var addr *net.UDPAddr
		var err error
		if IsServer {
			_, addr, err = udpConn.ReadFromUDP(data)
			if err != nil {
				panic(err)
			}
		} else {
			_, err = bufio.NewReader(udpConn).Read(data)
			if err != nil {
				panic(err)
			}
		}
		select {
		case incomings <- Incoming{data, addr}:
		default:
			log.Println("Error: queue is full, dropping message")
		}
	}
}

func Process() {
	var err error
	for {
		select {
		case incoming := <-incomings:
			var handler Handler
			switch incoming.data[0] {
			case JOIN:
				handler = &Join{}
			case MOVE:
				handler = &Move{}
			}
			handler.Unmarshal(incoming.data)
			if err = handler.Handle(incoming.addr); err != nil {
				log.Fatal(err)
			}
		default:
			// no data to process this frame
			return
		}
	}
}

func Send(data []byte, addr *net.UDPAddr) {
	if IsServer {
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

var curId int8 = 1

// Message type IDs
const (
	JOIN = iota
	MOVE
)

type Handler interface {
	Handle(addr *net.UDPAddr) error
	Marshal() []byte
	Unmarshal(b []byte)
}

type Join struct {
	ID  int8
	You bool
}

func (j *Join) Handle(addr *net.UDPAddr) error {
	player := NewPlayer()

	if IsServer {
		// player initialization (TODO set spawn point)
		player.addr = addr
		player.ID = curId
		curId++
		Lookup[addr.String()] = player.ID
		// tell the player their ID
		Send((&Join{ID: player.ID, You: true}).Marshal(), addr)
		// tell other players they joined
		for _, p := range Players {
			join := &Join{player.ID, false}
			Send(join.Marshal(), p.addr)
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

func (j *Join) Marshal() []byte {
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

type Move struct {
	ID                        int8
	AngularVelocity, Velocity float32
}

func (m *Move) Handle(addr *net.UDPAddr) error {
	var player *Player
	if IsServer {
		player = Players[Lookup[addr.String()]]
		if player == nil {
			log.Println("Player not found", addr.String(), Lookup[addr.String()])
			return nil
		}
		m.ID = player.ID
	} else {
		player = Players[m.ID]
		if player == nil {
			log.Println("Player not found with ID", m.ID)
			return nil
		}
	}

	player.Shape.Body.SetAngularVelocity(m.AngularVelocity)
	vx2 := math.Cos(float64(player.Shape.Body.Angle() * chipmunk.DegreeConst))
	vy2 := math.Sin(float64(player.Shape.Body.Angle() * chipmunk.DegreeConst))
	svx2, svy2 := m.Velocity*float32(vx2), m.Velocity*float32(vy2)
	player.Shape.Body.SetVelocity(svx2, svy2)

	if IsServer {
		theMove := m.Marshal()
		for _, p := range Players {
			Send(theMove, p.addr)
		}
	}

	return nil
}

func (m *Move) Marshal() []byte {
	buf := bytes.NewBuffer([]byte{MOVE, byte(m.ID)})
	err := binary.Write(buf, binary.LittleEndian, m.Velocity)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buf, binary.LittleEndian, m.AngularVelocity)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func (m *Move) Unmarshal(b []byte) {
	m.ID = int8(b[1])
	reader := bytes.NewReader(b[2:])
	err := binary.Read(reader, binary.LittleEndian, &m.Velocity)
	if err != nil {
		log.Println(err)
	}
	err = binary.Read(reader, binary.LittleEndian, &m.AngularVelocity)
	if err != nil {
		log.Println(err)
	}
}
