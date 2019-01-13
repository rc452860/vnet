package ssstream

import (
	"crypto/rand"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/rc452860/vnet/pool"
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
			buf:           make([]byte, MAX_PACKET_SIZE),
		}
		return sc, nil
	}
}

func (c *streamPacket) WriteTo(b []byte, addr net.Addr) (int, error) {
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
	encryper, err := c.NewEncrypter(c.key, iv)
	encryper.XORKeyStream(c.buf[ivLen:], b)
	_, err = c.PacketConn.WriteTo(c.buf[:ivLen+dataLen], addr)
	if err != nil {
		return 0, err
	}
	return dataLen, nil
}

func (c *streamPacket) ReadFrom(b []byte) (int, net.Addr, error) {
	n, addr, err := c.PacketConn.ReadFrom(b)
	if err != nil {
		return n, addr, err
	}
	ivLen := c.IVLen()

	if len(b) < ivLen {
		return n, addr, ErrShortPacket
	}

	decryptr, err := c.NewDecrypter(c.key, b[:ivLen])
	if err != nil {
		return n, addr, err
	}
	pool.GetUdpBuf()
	decryptr.XORKeyStream(b[ivLen:], b[ivLen:n])
	copy(b, b[ivLen:])
	return n - ivLen, addr, err
}

func (c *streamPacket) Close() error {
	pool.PutUdpBuf(c.buf)
	return c.PacketConn.Close()
}
