package gogame

import (
	"bufio"
	"bytes"
	"encoding"
	"encoding/binary"
	"log"
	"net"
)

const SimulatedNetworkLatencyMS = 100

// Message type IDs
const (
	JOIN = iota
	MOVE
	LOCATION
	PING
)

var ServerAddr *net.UDPAddr
var udpConn *net.UDPConn
var IsServer bool

type Incoming struct {
	handler Handler
	addr    *net.UDPAddr
}

type Outgoing struct {
	data []byte
	addr *net.UDPAddr
}

var incomings chan Incoming
var outgoings chan Outgoing

func init() {
	ServerAddr = &net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}
	incomings = make(chan Incoming, 100)
	outgoings = make(chan Outgoing, 100)
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

	if IsServer {
		go ProcessOutgoingServer()
	} else {
		go ProcessingOutgoingClient()
	}
}

func NetClose() error {
	return udpConn.Close()
}

// Recv runs in a goroutine and un-marshals incoming data, queuing it up for ProcessIncoming to handle
func Recv() {
	for {
		data := make([]byte, 2048)
		var addr *net.UDPAddr
		var err error
		if IsServer {
			_, addr, err = udpConn.ReadFromUDP(data)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			_, err = bufio.NewReader(udpConn).Read(data)
			if err != nil {
				log.Println(err)
				return
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
		//go func(){
		//	time.Sleep(SimulatedNetworkLatencyMS/2 * time.Millisecond)
		select {
		case incomings <- Incoming{handler, addr}:
		default:
			log.Println("Error: queue is full, dropping message")
		}
		//}()
	}
}

// ProcessIncoming handles all incoming messages that are queued
func ProcessIncoming() {
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

// Send queues up an outgoing byte array to be sent immediately so sending isn't blocking
func Send(data []byte, addr *net.UDPAddr) {
	//go func() {
	//	time.Sleep(SimulatedNetworkLatencyMS / 2 * time.Millisecond)
	outgoings <- Outgoing{data: data, addr: addr}
	//}()
}

func ProcessOutgoingServer() {
	var outgoing Outgoing
	for {
		outgoing = <-outgoings
		_, err := udpConn.WriteToUDP(outgoing.data, outgoing.addr)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func ProcessingOutgoingClient() {
	var outgoing Outgoing
	for {
		outgoing = <-outgoings
		_, err := udpConn.Write(outgoing.data)
		if err != nil {
			panic(err)
		}
	}
}

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
