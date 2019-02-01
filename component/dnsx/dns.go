package dnsx

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/rc452860/vnet/common/cache"
	"github.com/rc452860/vnet/common/log"
)

type DNS struct {
	dns1       string
	dns2       string
	IPV4Prefer bool
	*cache.Cache
}

const (
	IPV4 = "ipv4"
	IPV6 = "ipv6"
)

type DNSRecord struct {
	IP        net.IP
	IPType    string
	Domain    string
	expiredAt int64
}

func NewDNS(dns1, dns2 string) *DNS {
	return &DNS{
		dns1:  dns1,
		dns2:  dns2,
		Cache: cache.New(time.Second * 60),
	}
}

// NewDNSWithPrefer return DNS instance that take a dns reslove prefer
func NewDNSWithPrefer(dns1, dns2 string, ipv4Prefer bool) *DNS {
	dns := NewDNS(dns1, dns2)
	dns.IPV4Prefer = ipv4Prefer
	return dns
}

// MustReslove return ip or nil if not reslove dns it will return nil
func (d *DNS) MustReslove(domain string) net.IP {
	ip, err := d.Reslove(domain)
	if err != nil {
		log.Err(err)
		return nil
	}
	return ip
}

func (d *DNS) Reslove(domain string) (addr net.IP, err error) {
	key4 := fmt.Sprintf("%s,%s", IPV4, domain)
	key6 := fmt.Sprintf("%s,%s", IPV6, domain)

	if d.IPV4Prefer && d.Cache.Get(key4) != nil {
		if ipv4Reslove, ok := d.Cache.Get(key4).(*dns.A); ok {
			return ipv4Reslove.A, nil
		}
		if ipv6Reslove, ok := d.Cache.Get(key6).(*dns.A); ok {
			return ipv6Reslove.A, nil
		}
	}
	if !d.IPV4Prefer && d.Cache.Get(key6) != nil {
		if ipv6Reslove, ok := d.Cache.Get(key6).(*dns.A); ok {
			return ipv6Reslove.A, nil
		}
		if ipv4Reslove, ok := d.Cache.Get(key4).(*dns.A); ok {
			return ipv4Reslove.A, nil
		}
	}

	c := new(dns.Client)
	//TODO movie dns timeout to config
	// c.DialTimeout = time.Millisecond * 5000
	// c.ReadTimeout = time.Millisecond * 5000
	// c.WriteTimeout = time.Millisecond * 5000

	mipv4 := new(dns.Msg)
	mipv4.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	mipv6 := new(dns.Msg)
	mipv6.SetQuestion(dns.Fqdn(domain), dns.TypeAAAA)

	ch := make(chan *dns.Msg, 4)
	w := new(sync.WaitGroup)
	w.Add(4)
	go d.do(ch, w, c, mipv4, d.dns1)
	go d.do(ch, w, c, mipv4, d.dns2)
	go d.do(ch, w, c, mipv6, d.dns1)
	go d.do(ch, w, c, mipv6, d.dns2)
	w.Wait()
	close(ch)

	for response := range ch {
		if response.Rcode == dns.RcodeSuccess {
			for _, a := range response.Answer {
				switch t := a.(type) {
				case *dns.A:
					if d.Cache.Get(key4) == nil {
						d.Cache.Put(key4, t, time.Duration(t.Hdr.Ttl)*time.Second)
					}
				case *dns.AAAA:
					if d.Cache.Get(key6) == nil {
						d.Cache.Put(key6, t, time.Duration(t.Hdr.Ttl)*time.Second)
					}
				}

			}
		}
	}

	if d.IPV4Prefer {
		if ipv4Reslove, ok := d.Cache.Get(key4).(*dns.A); ok {
			return ipv4Reslove.A, nil
		}
		if ipv6Reslove, ok := d.Cache.Get(key6).(*dns.AAAA); ok {
			return ipv6Reslove.AAAA, nil
		}
	}
	if !d.IPV4Prefer {
		if ipv6Reslove, ok := d.Cache.Get(key6).(*dns.AAAA); ok {
			return ipv6Reslove.AAAA, nil
		}
		if ipv4Reslove, ok := d.Cache.Get(key4).(*dns.A); ok {
			return ipv4Reslove.A, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("dns can not reslove ip for domain: %s .", domain))
}

func (d *DNS) do(t chan *dns.Msg, wg *sync.WaitGroup, c *dns.Client, m *dns.Msg, addr string) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			log.Error("err:%v", err)
		}
	}()
	r, _, err := c.Exchange(m, addr)
	if err != nil {
		log.Error("%-v", err)
		return
	}
	t <- r
}
