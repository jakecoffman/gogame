package gogame

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"

	"github.com/vova616/chipmunk"
)

// Sent to server only: Move relays inputs related to movement
type Move struct {
	Turn, Throttle float32
}

func (m *Move) Handle(addr *net.UDPAddr) error {
	if !IsServer {
		log.Println("I shouldn't have gotten this")
		return nil
	}

	var player *Player = Players[Lookup[addr.String()]]
	if player == nil {
		log.Println("Player not found", addr.String(), Lookup[addr.String()])
		return nil
	}

	player.Shape.Body.SetAngularVelocity(m.Turn)
	vx2 := math.Cos(float64(player.Shape.Body.Angle() * chipmunk.DegreeConst))
	vy2 := math.Sin(float64(player.Shape.Body.Angle() * chipmunk.DegreeConst))
	svx2, svy2 := m.Throttle*float32(vx2), m.Throttle*float32(vy2)
	player.Shape.Body.SetVelocity(svx2, svy2)

	// Send immediate location update to everyone
	location := player.Location().Marshal()
	for _, p := range Players {
		Send(location, p.Addr)
	}

	return nil
}

func (m *Move) Marshal() []byte {
	buf := bytes.NewBuffer([]byte{MOVE})
	fields := []*float32{&m.Turn, &m.Throttle}
	var err error
	for _, field := range fields {
		err = binary.Write(buf, binary.LittleEndian, field)
		if err != nil {
			log.Fatal(err)
		}
	}
	return buf.Bytes()
}

func (m *Move) Unmarshal(b []byte) {
	reader := bytes.NewReader(b[1:])
	fields := []*float32{&m.Turn, &m.Throttle}
	for _, field := range fields {
		err := binary.Read(reader, binary.LittleEndian, field)
		if err != nil {
			log.Println(err)
		}
	}
}
