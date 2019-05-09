package network

import (
	"github.com/rs/xid"
	"net"
	"time"
)

func DialTcp(addr string) (req *Request, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return &Request{
		ISStream:    true,
		TCPConn:     conn,
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
		UDPConn:     conn,
		RequestID:   xid.New().String(),
		RequestTime: time.Now(),
	}, nil

}
