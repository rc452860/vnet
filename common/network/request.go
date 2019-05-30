package network

import (
	"github.com/rs/xid"
	"net"
	"time"
)

func NewRequestWithTCP(con net.Conn) *Request {
	request := new(Request)
	request.RequestID = xid.New().String()
	request.RequestTime = time.Now()
	request.ISStream = true
	request.Conn = con
	return request
}
func NewRequestWithUDP(con net.PacketConn) *Request {
	request := new(Request)
	request.RequestID = xid.New().String()
	request.RequestTime = time.Now()
	request.ISStream = false
	request.PacketConn = con
	return request
}

type Request struct {
	ISStream bool
	net.Conn
	net.PacketConn
	RequestID   string
	RequestTime time.Time
	Data        interface{}
}

func (r *Request) Close() error {
	if r.ISStream {
		return r.Conn.Close()
	} else {
		return r.PacketConn.Close()
	}
}

func (r *Request) LocalAddr() net.Addr {
	if r.ISStream {
		return r.Conn.LocalAddr()
	} else {
		return r.PacketConn.LocalAddr()
	}
}

func (r *Request) RemoteAddr() net.Addr {
	return r.Conn.RemoteAddr()
}

func (r *Request) SetDeadline(t time.Time) error {
	if r.ISStream {
		return r.Conn.SetDeadline(t)
	} else {
		return r.PacketConn.SetDeadline(t)
	}
}

func (r *Request) SetReadDeadline(t time.Time) error {
	if r.ISStream {
		return r.Conn.SetReadDeadline(t)
	} else {
		return r.PacketConn.SetReadDeadline(t)
	}
}

func (r *Request) SetWriteDeadline(t time.Time) error {
	if r.ISStream {
		return r.Conn.SetWriteDeadline(t)
	} else {
		return r.PacketConn.SetWriteDeadline(t)
	}
}

func (r *Request) SetKeepAlive(keepAlive bool) error{
	if tcpConn,ok := r.Conn.(*net.TCPConn);ok{
		return tcpConn.SetKeepAlive(keepAlive)
	}
	return nil
}
