package service

import (
	"github.com/rc452860/vnet/common/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// GetCpuUsage get cpu usage
func GetCPUUsage() int {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		log.Err(err)
		return 0
	}
	if len(percent) > 0 {
		return int(percent[0])
	} else {
		log.Error("get cpu usage fail")
		return 0
	}
}

//GetMemUsage get mem usage
func GetMemUsage() int {
	m, err := mem.VirtualMemory()
	if err != nil {
		log.Err(err)
		return 0
	}
	return int(m.UsedPercent)
}

//GetNetwork get Network traffic up and down
func GetNetwork() (up uint64, down uint64) {
	ni, err := net.IOCounters(false)
	if err != nil {
		log.Err(err)
		return 0, 0
	}
	if len(ni) > 0 {
		return ni[0].BytesSent, ni[0].BytesRecv
	} else {
		log.Error("can't get network tarffic")
		return 0, 0
	}
}

//GetDiskUsage get disk usage
func GetDiskUsage() int {
	d, err := disk.Usage("/")
	if err != nil {
		log.Err(err)
		return 0
	}
	return int(d.UsedPercent)
}
