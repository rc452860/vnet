package network

import (
	"github.com/rs/xid"
	"net"
	"time"
)

func DialTcp(addr string) (req *Request, err error) {
	conn, err := net.DialTimeout("tcp", addr,3 * time.Second)
	if err != nil {
		return nil, err
	}

	return &Request{
		ISStream:    true,
		Conn:     conn,
		RequestID:   xid.New().String(),
		RequestTime: time.Now(),
	}, nil
}

func DialUdp(addr string) (req *Request, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &Request{
		ISStream:    false,
		PacketConn:     conn,
		RequestID:   xid.New().String(),
		RequestTime: time.Now(),
	}, nil

}
