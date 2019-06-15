package service

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/rc452860/vnet/common/network"
	"github.com/rc452860/vnet/model"
	"github.com/rc452860/vnet/proxy"
	"github.com/rc452860/vnet/utils/addrx"
	"github.com/rc452860/vnet/utils/monitor"
	"github.com/sirupsen/logrus"
	"strings"
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
		uidPortTable:  make(map[int]int),
		portUidTable:  make(map[int]int),
		userInfo:      make(map[int]*model.UserInfo),
		UpTime:        time.Now(),
	}
}

type shadowsocksrService struct {
	Shadowsocksrs map[int]*proxy.ShadowsocksRProxy
	nodeInfo      *model.NodeInfo
	traffic       map[int]*model.UserTraffic
	trafficLock   *sync.Mutex
	online        map[int]*model.NodeOnline
	onlineLock    *sync.Mutex
	uidPortTable  map[int]int
	portUidTable  map[int]int
	userInfo      map[int]*model.UserInfo
	host          string
	UpTime        time.Time
}

// old python version shadowoscksr use port be protocol_param,
// so we must convert port to uid in web api version
func (s *shadowsocksrService) AddUIDPortConvertItem(uid, port int) {
	s.uidPortTable[uid] = port
	s.portUidTable[port] = uid
}

func (s *shadowsocksrService) DelUIDPortConvertItem(uid int) {
	port := s.uidPortTable[uid]
	delete(s.uidPortTable, uid)
	delete(s.portUidTable, port)
}

func (s *shadowsocksrService) UIDToPort(uid int) int {
	return s.uidPortTable[uid]
}

func (s *shadowsocksrService) PortToUid(port int) int {
	return s.portUidTable[port]
}

func (s *shadowsocksrService) Upload(port int, n int64) {
	s.trafficLock.Lock()
	uid := s.PortToUid(port)
	if s.traffic[uid] != nil {
		s.traffic[uid].Upload += n
	} else {
		traffic := new(model.UserTraffic)
		traffic.Upload += n
		traffic.Uid = uid
		s.traffic[uid] = traffic
	}
	s.trafficLock.Unlock()
}

func (s *shadowsocksrService) Download(port int, n int64) {
	s.trafficLock.Lock()
	uid := s.PortToUid(port)
	if s.traffic[uid] != nil {
		s.traffic[uid].Download += n
	} else {
		traffic := new(model.UserTraffic)
		traffic.Download += n
		traffic.Uid = uid
		s.traffic[uid] = traffic
	}
	s.trafficLock.Unlock()
}

func (s *shadowsocksrService) ReportTraffic() []*model.UserTraffic {
	s.trafficLock.Lock()
	reportData := s.traffic
	s.traffic = make(map[int]*model.UserTraffic)
	convertReportData := make([]*model.UserTraffic, 0, len(reportData))
	for _, value := range reportData {
		convertReportData = append(convertReportData, value)
	}
	s.trafficLock.Unlock()
	return convertReportData
}

func (s *shadowsocksrService) Online(port int, ip string) {
	s.onlineLock.Lock()
	uid := s.PortToUid(port)
	ip = addrx.SplitIpFromAddr(ip)
	if s.online[uid] == nil {
		nodeOnline := new(model.NodeOnline)
		nodeOnline.Uid = uid
		nodeOnline.IP = ip
		s.online[uid] = nodeOnline
	} else {
		if !strings.Contains(s.online[uid].IP, ip) {
			s.online[uid].IP = s.online[uid].IP + "," + ip
		}
	}
	s.onlineLock.Unlock()
}

func (s *shadowsocksrService) ReportOnline() []*model.NodeOnline {
	s.onlineLock.Lock()
	reportData := s.online
	convertReportData := make([]*model.NodeOnline, 0, len(reportData))
	for _, value := range reportData {
		convertReportData = append(convertReportData, value)
	}
	s.onlineLock.Unlock()
	return convertReportData
}

func (s *shadowsocksrService) ReportNodeStatus() model.NodeStatus {
	up, down := monitor.GetNetwork()
	return model.NodeStatus{
		CPU:    fmt.Sprintf("%v%%", monitor.GetCPUUsage()),
		MEM:    fmt.Sprintf("%v%%", monitor.GetMemUsage()),
		NET:    fmt.Sprintf("%v↑-%v↓", humanize.Bytes(up), humanize.Bytes(down)),
		DISK:   fmt.Sprintf("%v%%", monitor.GetDiskUsage()),
		UPTIME: int(time.Since(s.UpTime).Seconds()),
	}
}

func (s *shadowsocksrService) ShadowsocksRProxy(host string, port int, method, passwd, protocol, protocolParam, obfs, obfsParam string, args *proxy.ShadowsocksRArgs) {
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
	shadowsocksRProxy.OnlineReport = s
	shadowsocksRProxy.TrafficReport = s
	s.Shadowsocksrs[port] = shadowsocksRProxy
}

func (s *shadowsocksrService) SetNodeInfo(nodeInfo *model.NodeInfo) {
	s.nodeInfo = nodeInfo
}

func (s *shadowsocksrService) GetNodeInfo() *model.NodeInfo {
	return s.nodeInfo
}

func (s *shadowsocksrService) SetHost(host string) {
	s.host = host
}

func (s *shadowsocksrService) GetHost(host string) string {
	return s.host
}

func (s *shadowsocksrService) AddUser(user *model.UserInfo) {
	s.userInfo[user.Uid] = user
	if s.nodeInfo.Single == 1 {
		for _, server := range s.Shadowsocksrs {
			server.AddUser(user.Port, user.Passwd)
		}
	} else {
		if s.Shadowsocksrs[user.Port] != nil {
			if err := s.Shadowsocksrs[user.Port].Close(); err != nil {
				logrus.Error(err);
				return
			}
		}
		s.ShadowsocksRProxy(s.host,
			user.Port,
			s.nodeInfo.Method,
			user.Passwd,
			s.nodeInfo.Protocol,
			s.nodeInfo.ProtocolParam,
			s.nodeInfo.Obfs,
			s.nodeInfo.ObfsParam,
			&proxy.ShadowsocksRArgs{})
	}
}

func (s *shadowsocksrService) EditUser(user *model.UserInfo) {
	// TODO after change user profile it will be simultaneously exist old port and new port
	s.AddUser(user)
}

func (s *shadowsocksrService) DelUser(uid int) {
	port := s.UIDToPort(uid)
	if port == 0 {
		logrus.WithFields(logrus.Fields{
			"uid": uid,
		}).Info("uid is not exist")
		return
	}

	if s.nodeInfo.Single == 1 {
		for _, server := range s.Shadowsocksrs {
			server.DelUser(port)
			logrus.Infof("server %v del %v", server.Port, port)
		}
	} else {
		server := s.Shadowsocksrs[port]
		if server == nil {
			logrus.WithFields(logrus.Fields{
				"port": port,
			}).Info("port is not exist")
			return
		}

		if err := server.Close(); err != nil {
			logrus.Error(err)
			return
		}
	}
	logrus.WithFields(logrus.Fields{
		"uid": uid,
	}).Info("delete completed")
}
