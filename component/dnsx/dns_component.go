package dnsx

import (
	"github.com/rc452860/vnet/common/config"
	"github.com/rc452860/vnet/utils"
)

var (
	DNSComponent *DNS
)

func GetDNDComponent() *DNS {
	if DNSComponent == nil {
		InitDNSComponent()
	}
	return DNSComponent
}

func InitDNSComponent() {
	utils.Lock("DNSComponent")
	defer utils.UnLock("DNSComponent")
	conf := config.CurrentConfig().DNSOptions
	DNSComponent = NewDNSWithPrefer(conf.DNS1, conf.DNS2, conf.IPV4Prefer)
}
