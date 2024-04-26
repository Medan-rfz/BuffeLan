package multicast

import (
	"net"

	log "github.com/sirupsen/logrus"
)

const maxDatagramSize = 8192

type MulticastListener struct {
	conn  *net.UDPConn
	outCh chan OutputData
}

type OutputData struct {
	Src  *net.UDPAddr
	Data string
}

func NewMulticastListener(addr string) (*MulticastListener, error) {
	address, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, address)
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(maxDatagramSize)

	listener := &MulticastListener{
		conn:  conn,
		outCh: make(chan OutputData, 10),
	}

	go listener.listenStart()
	return listener, nil
}

func (c *MulticastListener) Listen() chan OutputData {
	return c.outCh
}

func (c *MulticastListener) Close() {
	c.conn.Close()
}

func (c *MulticastListener) listenStart() {
	defer close(c.outCh)

	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("ReadFromUDP failed:", err)
		}

		c.outCh <- msgHandler(src, numBytes, buffer)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) OutputData {
	return OutputData{
		Src:  src,
		Data: string(b[:n]),
	}
}
