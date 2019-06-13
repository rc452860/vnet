package ciphers

import (
	"fmt"
	"net"

	"github.com/rc452860/vnet/common/ciphers/aead"
	"github.com/rc452860/vnet/common/ciphers/stream"
)

//加密装饰
func CipherDecorate(password, method string, conn net.Conn) (net.Conn, error) {
	if method == "none" {
		return conn, nil
	}
	d := stream.GetStreamConnCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	d = aead.GetAEADConnCipher(method)
	if d != nil {
		return d(password, conn)
	}
	return nil, fmt.Errorf("[SS Cipher] not support : %s", method)
}

func CipherPacketDecorate(password, method string, conn net.PacketConn) (net.PacketConn, error) {
	if method == "none" {
		return conn, nil
	}
	d := stream.GetStreamPacketCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	d = aead.GetAEADPacketCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	return nil, fmt.Errorf("[SS Cipher] not support : %s", method)
}

func GetSupportCiphers() []string {
	stream := stream.GetStreamCiphers()
	list := make([]string, 0, 20)
	for k, _ := range stream {
		list = append(list, k)
	}
	aeas := aead.GetAEADCiphers()
	for k, _ := range aeas {
		list = append(list, k)
	}
	return list
}
