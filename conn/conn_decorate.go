package conn

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rc452860/vnet/pool"
	"github.com/rc452860/vnet/utils"
	"golang.org/x/time/rate"
)

var DefaultTimeOut = 10 * time.Second

func DialTcp(addr string) (IConn, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("can not connect addr %s with error %v", addr, err)
	}

	icon, err := DefaultDecorate(con, TCP)
	if err != nil {
		return nil, err
	}
	err = con.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		logging.Error("set tcp keepalive error: %v", err)
	}
	return icon, nil
}

func DefaultDecorate(c net.Conn, network string) (IConn, error) {
	id := utils.GetLongID()
	return &DefaultConn{
		Conn:    c,
		ID:      id,
		Network: network,
		context: context.Background(),
	}, nil
}

func DefaultDecorateForTls(c net.Conn, network string, id int64) (IConn, error) {
	return &DefaultConn{
		Conn:    c,
		ID:      id,
		Network: network,
		context: context.Background(),
	}, nil
}

type DefaultConn struct {
	net.Conn
	ID       int64
	RecordID int64
	Network  string
	context  context.Context
}

func (c *DefaultConn) GetID() int64 {
	return c.ID
}

func (c *DefaultConn) GetRecordID() int64 {
	return c.RecordID
}

func (c *DefaultConn) SetRecordID(id int64) {
	c.RecordID = id
}

func (c *DefaultConn) Flush() (int, error) {
	return 0, nil
}

func (c *DefaultConn) GetNetwork() string {
	return c.Network
}
func (c *DefaultConn) Context() context.Context {
	return c.context
}
func (c *DefaultConn) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *DefaultConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	return
}

func (c *DefaultConn) Write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}
func (c *DefaultConn) Close() error {
	return c.Conn.Close()
}

//超时装饰
func TimerDecorate(c IConn, rto, wto time.Duration) (IConn, error) {
	if rto == 0 {
		rto = DefaultTimeOut
	}
	if wto == 0 {
		wto = DefaultTimeOut
	}
	return &TimerConn{
		IConn:        c,
		ReadTimeOut:  rto,
		WriteTimeOut: wto,
	}, nil
}

type TimerConn struct {
	IConn
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

func (c *TimerConn) resetReadDeadline() {
	if c.ReadTimeOut > -1 {
		c.SetReadDeadline(time.Now().Add(c.ReadTimeOut))
	}
}

func (c *TimerConn) resetWriteDeadline() {
	if c.WriteTimeOut > -1 {
		c.SetWriteDeadline(time.Now().Add(c.WriteTimeOut))
	}
}

func (c *TimerConn) Read(b []byte) (n int, err error) {
	c.resetReadDeadline()
	n, err = c.IConn.Read(b)
	c.resetWriteDeadline()
	return
}

func (c *TimerConn) Write(b []byte) (n int, err error) {
	c.resetWriteDeadline()
	n, err = c.IConn.Write(b)
	c.resetReadDeadline()
	return
}

//缓冲装饰
func BufferDecorate(c IConn) (IConn, error) {
	return &BufferConn{
		IConn:  c,
		buffer: bytes.NewBuffer(pool.GetBuf()[:0]),
	}, nil
}

type BufferConn struct {
	IConn
	buffer *bytes.Buffer
}

func (c *BufferConn) Write(b []byte) (n int, err error) {
	return c.buffer.Write(b)
}

func (c *BufferConn) Flush() (n int, err error) {
	n, err = c.IConn.Write(c.buffer.Bytes())
	if err != nil {
		return
	}
	c.buffer.Reset()
	n, err = c.IConn.Flush()
	return
}

//实时写出
func RealTimeDecorate(c IConn) (IConn, error) {
	return &RealTimeFlush{
		IConn: c,
	}, nil
}

type RealTimeFlush struct {
	IConn
}

func (r *RealTimeFlush) Write(b []byte) (n int, err error) {
	n, err = r.IConn.Write(b)
	if err != nil {
		return
	}
	_, err = r.IConn.Flush()
	return
}

//导出装饰器
type TrafficHandle func(IConn, int)

func TrafficDecorate(c IConn, upload, download TrafficHandle) (IConn, error) {
	return &Traffic{
		IConn:    c,
		Upload:   upload,
		Download: download,
	}, nil
}

type Traffic struct {
	IConn
	Upload   TrafficHandle
	Download TrafficHandle
}

func (t *Traffic) Read(b []byte) (n int, err error) {
	n, err = t.IConn.Read(b)
	if t.Download != nil {
		t.Download(t.IConn, n)
	}
	return
}

func (t *Traffic) Write(b []byte) (n int, err error) {
	n, err = t.IConn.Write(b)
	if t.Upload != nil {
		t.Upload(t.IConn, n)
	}
	return
}

// traffic limit decorate
type TrafficLimit struct {
	IConn
	ReadLimit  *rate.Limiter
	WriteLimit *rate.Limiter
}

func TrafficLimitDecorate(con IConn, read, write *rate.Limiter) (IConn, error) {
	return &TrafficLimit{
		IConn:      con,
		ReadLimit:  read,
		WriteLimit: write,
	}, nil
}

func (c *TrafficLimit) Read(b []byte) (n int, err error) {
	if c.ReadLimit == nil {
		return c.IConn.Read(b)
	}
	n, err = c.IConn.Read(b)
	if err != nil {
		return n, err
	}
	if err = c.ReadLimit.WaitN(context.Background(), n); err != nil {
		return n, err
	}
	return n, nil
}

func (c *TrafficLimit) Write(b []byte) (n int, err error) {
	if c.WriteLimit == nil {
		return c.IConn.Write(b)
	}
	n, err = c.IConn.Write(b)
	if err != nil {
		return n, err
	}
	if err = c.WriteLimit.WaitN(context.Background(), n); err != nil {
		return n, err
	}
	return n, nil
}

// for blow code is about udp communicateion
// is decorate packet con to provide traffic record and traffic limit etc ...
type PacketTrafficHandle func(laddr string, raddr string, n int)

type PacketTrafficConn struct {
	net.PacketConn
	upload   PacketTrafficHandle
	download PacketTrafficHandle
}

func PacketTrafficConnDecorate(con net.PacketConn, upload, download PacketTrafficHandle) net.PacketConn {
	return &PacketTrafficConn{
		PacketConn: con,
		upload:     upload,
		download:   download,
	}
}

func (this *PacketTrafficConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, addr, err = this.PacketConn.ReadFrom(p)
	if err != nil {
		return n, addr, err
	}
	if this.download != nil {
		this.download(this.PacketConn.LocalAddr().String(), addr.String(), n)
	}
	return n, addr, err
}

func (this *PacketTrafficConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	n, err = this.PacketConn.WriteTo(p, addr)
	if this.upload != nil {
		this.upload(this.PacketConn.LocalAddr().String(), addr.String(), n)
	}
	return n, err
}
