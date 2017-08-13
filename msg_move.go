package gogame

import (
	"bytes"
	"log"
	"math"
	"net"
	"github.com/jakecoffman/physics"
)

// Sent to server only: Move relays inputs related to movement
type Move struct {
	Turn, Throttle float64
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

	player.Shape.Body().SetAngularVelocity(m.Turn)
	vx2 := math.Cos(float64(player.Shape.Body().Angle() * physics.DegreeConst))
	vy2 := math.Sin(float64(player.Shape.Body().Angle() * physics.DegreeConst))
	svx2, svy2 := m.Throttle*vx2, m.Throttle*vy2
	player.Shape.Body().SetVelocity(svx2, svy2)

	// Send immediate location update to everyone
	location, err := player.Location().MarshalBinary()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, p := range Players {
		Send(location, p.Addr)
	}

	return nil
}

func (m *Move) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{MOVE})
	fields := []interface{}{&m.Turn, &m.Throttle}
	return Marshal(fields, buf)
}

func (m *Move) UnmarshalBinary(b []byte) error {
	reader := bytes.NewReader(b[1:])
	fields := []interface{}{&m.Turn, &m.Throttle}
	return Unmarshal(fields, reader)
}
