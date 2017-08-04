package gogame

import (
	"bytes"
	"log"
	"net"
	"time"
)

type Ping struct {
	Sent time.Time
}

func (p *Ping) Handle(addr *net.UDPAddr) error {
	if IsServer {
		// got a ping
		player := Players[Lookup[addr.String()]]
		log.Println("Player", player.ID, "ping", time.Since(p.Sent)*2)
		bin, err := p.MarshalBinary()
		if err != nil {
			log.Println(err)
			return err
		}
		Send(bin, addr)
	} else {
		// got a pong
		log.Println("Ping", time.Since(p.Sent))
	}
	return nil
}

func (p *Ping) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{PING})
	return Marshal([]interface{}{p.Sent.UnixNano()}, buf)
}

func (p *Ping) UnmarshalBinary(b []byte) error {
	var nsecs int64
	err := Unmarshal([]interface{}{&nsecs}, bytes.NewReader(b[1:]))
	if err != nil {
		log.Println(err)
		return err
	}
	p.Sent = time.Unix(0, nsecs)
	return nil
}

var done chan struct{}

func PingRegularly() {
	done = make(chan struct{})
	tick := time.Tick(3 * time.Second)
	for {
		select {
		case <-done:
			close(done)
			return
		case <-tick:
			ping := &Ping{Sent: time.Now()}
			bin, err := ping.MarshalBinary()
			if err != nil {
				log.Println(err)
				return
			}
			Send(bin, ServerAddr)
		}
	}
}
