package gogame

import (
	"bytes"
	"log"
	"net"

	"github.com/jakecoffman/physics"
)

// message sent to clients: update location information
type Location struct {
	ID                     int8
	X, Y                   float64
	Vx, Vy                 float64
	Angle, AngularVelocity float64
}

func (l *Location) Handle(addr *net.UDPAddr) error {
	if IsServer {
		log.Println("I shouldn't have gotten this")
		return nil
	}

	player := Players[l.ID]
	if player == nil {
		log.Println("Player with ID", l.ID, "not found")
		return nil
	}
	// TODO: check if the change is insignificant and ignore it if that's the case
	body := player.Shape.Body()
	body.SetPosition(&physics.Vector{l.X, l.Y})
	body.SetVelocity(l.Vx, l.Vy)
	body.SetAngle(l.Angle)
	body.SetAngularVelocity(l.AngularVelocity)
	return nil
}

func (l *Location) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{LOCATION, byte(l.ID)})
	fields := []interface{}{&l.X, &l.Y, &l.Vx, &l.Vy, &l.Angle, &l.AngularVelocity}
	return Marshal(fields, buf)
}

func (l *Location) UnmarshalBinary(b []byte) error {
	l.ID = int8(b[1])
	reader := bytes.NewReader(b[2:])
	fields := []interface{}{&l.X, &l.Y, &l.Vx, &l.Vy, &l.Angle, &l.AngularVelocity}
	return Unmarshal(fields, reader)
}
