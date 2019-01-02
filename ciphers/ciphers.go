package ciphers

import (
	"fmt"
	"net"

	"github.com/rc452860/vnet/ciphers/ssaead"
	"github.com/rc452860/vnet/ciphers/ssstream"
	connect "github.com/rc452860/vnet/conn"
)

type ConnDecorate func(password string, conn connect.IConn) (connect.IConn, error)

//加密装饰
func CipherDecorate(password, method string, conn connect.IConn) (connect.IConn, error) {
	d := ssstream.GetStreamConnCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	d = ssaead.GetAEADConnCipher(method)
	if d != nil {
		return d(password, conn)
	}
	return nil, fmt.Errorf("[SS Cipher] not support : %s", method)
}

func CipherPacketDecorate(password, method string, conn net.PacketConn) (net.PacketConn, error) {
	d := ssstream.GetStreamPacketCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	d = ssaead.GetAEADPacketCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	return nil, fmt.Errorf("[SS Cipher] not support : %s", method)
}
