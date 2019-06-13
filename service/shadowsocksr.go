package service

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/rc452860/vnet/common/network"
	"github.com/rc452860/vnet/model"
	"github.com/rc452860/vnet/proxy"
	"github.com/rc452860/vnet/utils/monitor"
	"sync"
	"time"
)

var ShadowsocksrServiceInstance *shadowsocksrService

func init() {
	ShadowsocksrServiceInstance = &shadowsocksrService{
		Shadowsocksrs: make(map[int]*proxy.ShadowsocksRProxy),
		traffic:       make(map[int]*model.UserTraffic),
		trafficLock:   new(sync.Mutex),
		online:        make(map[int]*model.NodeOnline),
		onlineLock:    new(sync.Mutex),
		uidPortTable: make(map[int]int),
		portUidTable: make(map[int]int),
	}
}

type shadowsocksrService struct {
	Shadowsocksrs map[int]*proxy.ShadowsocksRProxy
	traffic       map[int]*model.UserTraffic
	trafficLock   *sync.Mutex
	online        map[int]*model.NodeOnline
	onlineLock    *sync.Mutex
	uidPortTable map[int]int
	portUidTable map[int]int
}

// old python version shadowoscksr use port be protocol_param,
// so we must convert port to uid in web api version
func (ssrService *shadowsocksrService) AddUIDPortConvertItem(uid,port int){
	ssrService.uidPortTable[uid] = port
	ssrService.portUidTable[port] = uid
}

func (ssrService *shadowsocksrService) DelUIDPortConvertItem(uid int){
	port := ssrService.uidPortTable[uid]
	delete(ssrService.uidPortTable,uid)
	delete(ssrService.portUidTable,port)
}

func (ssrService *shadowsocksrService) UIDToPort(uid int) int{
	return ssrService.uidPortTable[uid]
}

func (ssrService *shadowsocksrService) PortToUid(port int)int{
	return ssrService.portUidTable[port]
}

func (ssrService *shadowsocksrService) Upload(port int, n int64) {
	ssrService.trafficLock.Lock()
	uid := ssrService.PortToUid(port)
	if ssrService.traffic[uid] != nil {
		ssrService.traffic[uid].Upload += n
	} else {
		traffic := new(model.UserTraffic)
		traffic.Upload += n
		traffic.Uid = uid
		ssrService.traffic[uid] = traffic
	}
	ssrService.trafficLock.Unlock()
}

func (ssrService *shadowsocksrService) Download(port int, n int64) {
	ssrService.trafficLock.Lock()
	uid := ssrService.PortToUid(port)
	if ssrService.traffic[uid] != nil {
		ssrService.traffic[uid].Download += n
	} else {
		traffic := new(model.UserTraffic)
		traffic.Download += n
		traffic.Uid = uid
		ssrService.traffic[uid] = traffic
	}
	ssrService.trafficLock.Unlock()
}

func (ssrService *shadowsocksrService) ReportTraffic() []*model.UserTraffic {
	ssrService.trafficLock.Lock()
	reportData := ssrService.traffic
	ssrService.traffic = make(map[int]*model.UserTraffic)
	convertReportData := make([]*model.UserTraffic, 0, len(reportData))
	for _, value := range reportData {
		convertReportData = append(convertReportData, value)
	}
	ssrService.trafficLock.Unlock()
	return convertReportData
}

func (ssrService *shadowsocksrService) Online(port int, ip string) {
	ssrService.onlineLock.Lock()
	uid := ssrService.PortToUid(port)
	if ssrService.online[uid] == nil {
		nodeOnline := new(model.NodeOnline)
		nodeOnline.Uid = uid
		nodeOnline.IP = ip
		ssrService.online[uid] = nodeOnline
	} else {
		ssrService.online[uid].IP = ssrService.online[uid].IP + "," + ip
	}
	ssrService.onlineLock.Unlock()
}

func (ssrService *shadowsocksrService) ReportOnline() []*model.NodeOnline {
	ssrService.onlineLock.Lock()
	reportData := ssrService.online
	convertReportData := make([]*model.NodeOnline,0,len(reportData))
	for _,value := range reportData{
		convertReportData = append(convertReportData,value)
	}
	ssrService.onlineLock.Unlock()
	return convertReportData
}

func (ssrService *shadowsocksrService) ReportNodeStatus() model.NodeStatus {
	up, down := monitor.GetNetwork()
	return model.NodeStatus{
		CPU:  fmt.Sprintf("%v%%", monitor.GetCPUUsage()),
		MEM:  fmt.Sprintf("%v%%", monitor.GetMemUsage()),
		NET:  fmt.Sprintf("%v↑-%v↓", humanize.Bytes(up), humanize.Bytes(down)),
		DISK: fmt.Sprintf("%v%%", monitor.GetDiskUsage()),
	}
}

func (ssrService *shadowsocksrService) ShadowsocksRProxy(host string, port int, method, passwd, protocol, protocolParam, obfs, obfsParam string, args *proxy.ShadowsocksRArgs) {
	shadowsocksRProxy := new(proxy.ShadowsocksRProxy)
	shadowsocksRProxy.Host = host
	shadowsocksRProxy.Port = port
	shadowsocksRProxy.Method = method
	shadowsocksRProxy.Password = passwd
	shadowsocksRProxy.Protocol = protocol
	shadowsocksRProxy.ProtocolParam = protocolParam
	shadowsocksRProxy.Obfs = obfs
	shadowsocksRProxy.ObfsParam = obfsParam
	shadowsocksRProxy.ShadowsocksRArgs = args
	shadowsocksRProxy.Listener = network.NewListener(fmt.Sprintf("%s:%v", host, port), 3*time.Second)
	shadowsocksRProxy.OnlineReport = ssrService
	shadowsocksRProxy.TrafficReport = ssrService
	ssrService.Shadowsocksrs[port] = shadowsocksRProxy
}
