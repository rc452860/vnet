package network

import (
	"github.com/rs/xid"
	"net"
	"time"
)

func NewRequest(isstream bool, con net.Conn) *Request {
	request := new(Request)
	request.RequestID = xid.New().String()
	request.RequestTime = time.Now()
	request.ISStream = isstream
	if isstream {
		request.TCPConn = con.(*net.TCPConn)
	} else {
		request.UDPConn = con.(*net.UDPConn)
	}
	return request
}

type Request struct {
	ISStream bool
	*net.TCPConn
	*net.UDPConn
	RequestID   string
	RequestTime time.Time
	Data        interface{}
}

func (r *Request) Close() error {
	if r.ISStream {
		return r.TCPConn.Close()
	} else {
		return r.UDPConn.Close()
	}
}

func (r *Request) LocalAddr() net.Addr {
	if r.ISStream {
		return r.TCPConn.LocalAddr()
	} else {
		return r.UDPConn.LocalAddr()
	}
}

func (r *Request) RemoteAddr() net.Addr {
	if r.ISStream {
		return r.TCPConn.RemoteAddr()
	} else {
		return r.UDPConn.RemoteAddr()
	}
}

func (r *Request) SetDeadline(t time.Time) error {
	if r.ISStream {
		return r.TCPConn.SetDeadline(t)
	} else {
		return r.UDPConn.SetDeadline(t)
	}
}

func (r *Request) SetReadDeadline(t time.Time) error {
	if r.ISStream {
		return r.TCPConn.SetReadDeadline(t)
	} else {
		return r.UDPConn.SetReadDeadline(t)
	}
}

func (r *Request) SetWriteDeadline(t time.Time) error {
	if r.ISStream {
		return r.TCPConn.SetWriteDeadline(t)
	} else {
		return r.UDPConn.SetWriteDeadline(t)
	}
}

func (r *Request) Read(b []byte) (n int, err error){
	if r.ISStream {
		return r.TCPConn.Read(b)
	} else {
		return r.UDPConn.Read(b)
	}
}

func (r *Request) Write(b []byte) (n int, err error){
	if r.ISStream {
		return r.TCPConn.Write(b)
	} else {
		return r.UDPConn.Write(b)
	}
}
