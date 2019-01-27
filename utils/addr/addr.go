package addr

import (
	"net"

	"github.com/rc452860/vnet/common/log"
)

func GetIPFromAddr(addr net.Addr) string {
	switch addr.(type) {
	case *net.TCPAddr:
		tcpAddr := addr.(*net.TCPAddr)
		return tcpAddr.IP.String()
	case *net.UDPAddr:
		udpAddr := addr.(*net.UDPAddr)
		return udpAddr.IP.String()
	case nil:
		return ""
	default:
		return ""
	}
}

func GetPortFromAddr(addr net.Addr) int {
	switch addr.(type) {
	case *net.TCPAddr:
		tcpAddr := addr.(*net.TCPAddr)
		return tcpAddr.Port
	case *net.UDPAddr:
		udpAddr := addr.(*net.UDPAddr)
		return udpAddr.Port
	case nil:
		return 0
	default:
		return 0
	}
}

func GetNetworkFromAddr(addr net.Addr) string {
	return addr.Network()
}

func ParseAddrFromString(network, addr string) net.Addr {
	var addrConvert net.Addr
	var err error
	switch network {
	case "tcp", "tcp4", "tcp6":
		addrConvert, err = net.ResolveTCPAddr(network, addr)
	case "udp", "udp4", "udp6":
		addrConvert, err = net.ResolveUDPAddr(network, addr)
	}
	if err != nil {
		log.Err(err)
		return nil
	}
	return addrConvert
}
