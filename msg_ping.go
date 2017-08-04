package gogame

import (
	"bytes"
	"log"
	"net"
	"sync"
	"time"
)

var LastPing *lastPing

func init() {
	LastPing = &lastPing{}
}

type lastPing struct {
	sync.RWMutex
	Duration time.Duration
}

func (l *lastPing) Set(d time.Duration) {
	l.Lock()
	l.Duration = d
	l.Unlock()
}

func (l *lastPing) Get() time.Duration {
	l.RLock()
	defer l.RUnlock()
	return l.Duration
}

type Ping struct {
	Sent time.Time
}

func (p *Ping) Handle(addr *net.UDPAddr) error {
	if IsServer {
		// got a ping
		// TODO store pings for players to show on server status
		//player := Players[Lookup[addr.String()]]
		bin, err := p.MarshalBinary()
		if err != nil {
			log.Println(err)
			return err
		}
		Send(bin, addr)
	} else {
		// got a pong
		d := time.Since(p.Sent)
		LastPing.Set(d)
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
	PingNow()
	done = make(chan struct{})
	tick := time.Tick(3 * time.Second)
	for {
		select {
		case <-done:
			close(done)
			return
		case <-tick:
			PingNow()
		}
	}
}

func PingNow() {
	ping := &Ping{Sent: time.Now()}
	bin, err := ping.MarshalBinary()
	if err != nil {
		log.Println(err)
		return
	}
	Send(bin, ServerAddr)
}
