package multicast

import "net"

type MulticastSender struct {
	conn *net.UDPConn
}

func NewMulticastSender(addr string) (*MulticastSender, error) {
	address, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, address)
	if err != nil {
		return nil, err
	}

	return &MulticastSender{
		conn: conn,
	}, nil
}

func (s *MulticastSender) Send(data string) {
	s.conn.Write([]byte(data))
}

func (s *MulticastSender) Close() {
	s.conn.Close()
}
