package gogame

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"encoding"
)

var ServerAddr *net.UDPAddr
var udpConn *net.UDPConn
var IsServer bool

type Incoming struct {
	handler Handler
	addr    *net.UDPAddr
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
		var handler Handler
		switch data[0] {
		case JOIN:
			handler = &Join{}
		case MOVE:
			handler = &Move{}
		case LOCATION:
			handler = &Location{}
		case PING:
			handler = &Ping{}
			// just handle the ping here immediately outside of the game loop
			err = handler.UnmarshalBinary(data)
			if err != nil {
				log.Println(err)
				continue
			}
			handler.Handle(addr)
			continue
		default:
			log.Println("Unkown message type", data[0])
			continue
		}
		err = handler.UnmarshalBinary(data)
		if err != nil {
			log.Println(err)
			continue
		}
		select {
		case incomings <- Incoming{handler, addr}:
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
			if err = incoming.handler.Handle(incoming.addr); err != nil {
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

// Message type IDs
const (
	JOIN = iota
	MOVE
	LOCATION
	PING
)

type Handler interface {
	Handle(addr *net.UDPAddr) error
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func Marshal(fields []interface{}, buf *bytes.Buffer) ([]byte, error) {
	var err error
	for _, field := range fields {
		err = binary.Write(buf, binary.LittleEndian, field)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func Unmarshal(fields []interface{}, reader *bytes.Reader) error {
	for _, field := range fields {
		err := binary.Read(reader, binary.LittleEndian, field)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
