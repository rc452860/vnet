// Package socks implements essential parts of SOCKS protocol.
package socks

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

// UDPEnabled is the toggle for UDP support
var UDPEnabled = true

// SOCKS request commands as defined in RFC 1928 section 4.
const (
	CmdConnect      = 1
	CmdBind         = 2
	CmdUDPAssociate = 3
)

// SOCKS address types as defined in RFC 1928 section 5.
const (
	AtypIPv4       = 1
	AtypDomainName = 3
	AtypIPv6       = 4
)

// Error represents a SOCKS error
type Error byte

func (err Error) Error() string {
	return "SOCKS error: " + strconv.Itoa(int(err))
}

// SOCKS errors as defined in RFC 1928 section 6.
const (
	ErrGeneralFailure       = Error(1)
	ErrConnectionNotAllowed = Error(2)
	ErrNetworkUnreachable   = Error(3)
	ErrHostUnreachable      = Error(4)
	ErrConnectionRefused    = Error(5)
	ErrTTLExpired           = Error(6)
	ErrCommandNotSupported  = Error(7)
	ErrAddressNotSupported  = Error(8)
	InfoUDPAssociate        = Error(9)
)

// MaxAddrLen is the maximum size of SOCKS address in bytes.
const MaxAddrLen = 1 + 1 + 255 + 2

// Socks5Addr represents a SOCKS address as defined in RFC 1928 section 5.
type Socks5Addr struct {
	Raw     []byte
	AType   int
	Address string
	Port    int
}

func NewSocks5Addr(raw []byte, atype int) *Socks5Addr {
	addr := &Socks5Addr{
		Raw:   raw,
		AType: atype,
	}
	addr.process()
	return addr
}

func (s *Socks5Addr) GetAddress() string {
	return s.Address
}

func (s *Socks5Addr) GetPort() int {
	return s.Port
}

func (s *Socks5Addr) GetAType() int {
	return s.AType
}

func (s *Socks5Addr) process() {
	switch s.AType { // address type
	case AtypDomainName:
		s.Address = string(s.Raw[2 : 2+int(s.Raw[1])])
		s.Port = (int(s.Raw[2+int(s.Raw[1])]) << 8) | int(s.Raw[2+int(s.Raw[1])+1])
	case AtypIPv4:
		s.Address = net.IP(s.Raw[1 : 1+net.IPv4len]).String()
		s.Port = (int(s.Raw[1+net.IPv4len]) << 8) | int(s.Raw[1+net.IPv4len+1])
	case AtypIPv6:
		s.Address = net.IP(s.Raw[1 : 1+net.IPv6len]).String()
		s.Port = (int(s.Raw[1+net.IPv6len]) << 8) | int(s.Raw[1+net.IPv6len+1])
	}
}

func (s *Socks5Addr) String() string {
	return fmt.Sprintf("%s:%v", s.Address, s.Port)
}

func readAddr(r io.Reader, b []byte) (*Socks5Addr, error) {
	if len(b) < MaxAddrLen {
		return nil, io.ErrShortBuffer
	}
	_, err := io.ReadFull(r, b[:1]) // read 1st byte for address type
	if err != nil {
		return nil, err
	}

	switch b[0] {
	case AtypDomainName:
		_, err = io.ReadFull(r, b[1:2]) // read 2nd byte for domain length
		if err != nil {
			return nil, err
		}
		_, err = io.ReadFull(r, b[2:2+int(b[1])+2])
		return NewSocks5Addr(b[:1+1+int(b[1])+2], AtypDomainName), err
	case AtypIPv4:
		_, err = io.ReadFull(r, b[1:1+net.IPv4len+2])
		return NewSocks5Addr(b[:1+net.IPv4len+2], AtypIPv4), err
	case AtypIPv6:
		_, err = io.ReadFull(r, b[1:1+net.IPv6len+2])
		return NewSocks5Addr(b[:1+net.IPv6len+2], AtypIPv6), err
	}

	return nil, ErrAddressNotSupported
}

// ReadAddr reads just enough bytes from r to get a valid Addr.
func ReadAddr(r io.Reader) (*Socks5Addr, error) {
	return readAddr(r, make([]byte, MaxAddrLen))
}

// SplitAddr slices a SOCKS address from beginning of b. Returns nil if failed.
func SplitAddr(b []byte) *Socks5Addr {
	addrLen := 1
	if len(b) < addrLen {
		return nil
	}

	var atype int
	switch b[0] {
	case AtypDomainName:
		if len(b) < 2 {
			return nil
		}
		addrLen = 1 + 1 + int(b[1]) + 2
		atype = AtypDomainName
	case AtypIPv4:
		addrLen = 1 + net.IPv4len + 2
		atype = AtypIPv4
	case AtypIPv6:
		addrLen = 1 + net.IPv6len + 2
		atype = AtypIPv6
	default:
		return nil

	}

	if len(b) < addrLen {
		return nil
	}

	return NewSocks5Addr(b[:addrLen], atype)
}

// ParseAddr parses the address in string s. Returns nil if failed.
func ParseAddr(s string) *Socks5Addr {
	var (
		addr  []byte
		aType int
	)
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			addr = make([]byte, 1+net.IPv4len+2)
			addr[0] = AtypIPv4
			copy(addr[1:], ip4)
			aType = AtypIPv4
		} else {
			addr = make([]byte, 1+net.IPv6len+2)
			addr[0] = AtypIPv6
			copy(addr[1:], ip)
			aType = AtypIPv6
		}
	} else {
		if len(host) > 255 {
			return nil
		}
		addr = make([]byte, 1+1+len(host)+2)
		addr[0] = AtypDomainName
		addr[1] = byte(len(host))
		copy(addr[2:], host)
		aType = AtypDomainName
	}

	portnum, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	addr[len(addr)-2], addr[len(addr)-1] = byte(portnum>>8), byte(portnum)

	return NewSocks5Addr(addr, aType)
}
