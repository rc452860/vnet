package stream

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/rc452860/vnet/common/pool"
)

var ErrShortPacket = errors.New("short packet")

const MAX_PACKET_SIZE = 64 * 1204

type streamPacket struct {
	net.PacketConn
	IStreamCipher
	sync.Mutex
	key []byte
	buf []byte
}

func GetStreamPacketCiphers(method string) func(string, net.PacketConn) (net.PacketConn, error) {
	c, ok := streamCiphers[method]
	if !ok {
		return nil
	}
	return func(password string, packet net.PacketConn) (net.PacketConn, error) {
		iv := make([]byte, c.IVLen())
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		sc := &streamPacket{
			PacketConn:    packet,
			IStreamCipher: c,
			key:           evpBytesToKey(password, c.KeyLen()),
			buf:           pool.GetBuf(),
		}
		return sc, nil
	}
}

func (c *streamPacket) GetKey() []byte {
	return c.key
}

func (c *streamPacket) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	c.Lock()
	defer c.Unlock()
	ivLen := c.IVLen()
	dataLen := len(b)
	if MAX_PACKET_SIZE < ivLen+dataLen {
		return 0, io.ErrShortBuffer
	}
	iv := c.buf[:ivLen]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}
	encryper, err := c.NewStream(c.key, iv, 0)
	encryper.XORKeyStream(c.buf[ivLen:], b)
	_, err = c.PacketConn.WriteTo(c.buf[:ivLen+dataLen], addr)
	if err != nil {
		return 0, err
	}
	return dataLen, nil
}

func (c *streamPacket) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.WithStack(errors.New(fmt.Sprintf("%v", e)))
		}
	}()
	n, addr, err = c.PacketConn.ReadFrom(b)

	if err != nil {
		return n, addr, err
	}
	ivLen := c.IVLen()

	if len(b) > cap(c.buf) {
		return 0, nil, errors.WithStack(io.ErrShortBuffer)
	}

	if n < ivLen || len(b) < ivLen {
		return n, addr, ErrShortPacket
	}

	decryptr, err := c.NewStream(c.key, b[:ivLen], 0)
	if err != nil {
		return n, addr, err
	}

	decryptr.XORKeyStream(b[ivLen:], b[ivLen:ivLen+n])
	copy(b, b[ivLen:])
	return n - ivLen, addr, err
}

func (c *streamPacket) Close() error {
	pool.PutBuf(c.buf)
	return c.PacketConn.Close()
}
