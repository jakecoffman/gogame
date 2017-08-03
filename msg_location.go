package gogame

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"github.com/vova616/chipmunk/vect"
)

// message sent to clients: update location information
type Location struct {
	ID                     int8
	X, Y                   float32
	Vx, Vy                 float32
	Angle, AngularVelocity float32
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
	body := player.Shape.Body
	body.SetPosition(vect.Vect{vect.Float(l.X), vect.Float(l.Y)})
	body.SetVelocity(l.Vx, l.Vy)
	body.SetAngle(vect.Float(l.Angle))
	body.SetAngularVelocity(l.AngularVelocity)
	return nil
}

func (l *Location) Marshal() []byte {
	buf := bytes.NewBuffer([]byte{LOCATION, byte(l.ID)})
	fields := []*float32{&l.X, &l.Y, &l.Vx, &l.Vy, &l.Angle, &l.AngularVelocity}
	var err error
	for _, field := range fields {
		err = binary.Write(buf, binary.LittleEndian, field)
		if err != nil {
			log.Fatal(err)
		}
	}
	return buf.Bytes()
}

func (l *Location) Unmarshal(b []byte) {
	l.ID = int8(b[1])
	reader := bytes.NewReader(b[2:])
	fields := []*float32{&l.X, &l.Y, &l.Vx, &l.Vy, &l.Angle, &l.AngularVelocity}
	for _, field := range fields {
		err := binary.Read(reader, binary.LittleEndian, field)
		if err != nil {
			log.Println(err)
		}
	}
}
