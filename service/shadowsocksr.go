package service

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/rc452860/vnet/common/network"
	"github.com/rc452860/vnet/model"
	"github.com/rc452860/vnet/proxy"
	"github.com/rc452860/vnet/utils/addrx"
	"github.com/rc452860/vnet/utils/monitor"
	"github.com/sirupsen/logrus"
	"strconv"
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
		userTable:     make(map[int]*model.UserInfo),
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
	userTable     map[int]*model.UserInfo
	host          string
	UpTime        time.Time
}

func (s *shadowsocksrService) UIDToPort(uid int) int {
	if user := s.userTable[uid]; user != nil {
		return user.Port
	}
	return 0
}

func (s *shadowsocksrService) PortToUid(port int) int {
	for _, value := range s.userTable {
		if value.Port == port {
			return value.Uid
		}
	}
	return 0
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
	for key, value := range reportData {
		if value.Download + value.Upload < 50 * 1024{
			continue
		}
		convertReportData = append(convertReportData, value)
		delete(s.traffic,key)
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

func (s *shadowsocksrService) ShadowsocksRProxy(host string, port int, method, passwd, protocol, protocolParam, obfs, obfsParam string,single int, args *proxy.ShadowsocksRArgs) *proxy.ShadowsocksRProxy {
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
	shadowsocksRProxy.Single = single
	s.Shadowsocksrs[port] = shadowsocksRProxy
	return shadowsocksRProxy
}

func (s *shadowsocksrService) SetNodeInfo(nodeInfo *model.NodeInfo) {
	s.nodeInfo = nodeInfo
	if s.nodeInfo.Single == 1 {
		portStrArray := strings.Split(nodeInfo.Port, ",")
		ports := []int{}
		for _, item := range portStrArray {
			convertPort, err := strconv.Atoi(item)
			if err != nil {
				panic(fmt.Sprintf("port format error: %s", nodeInfo.Port))
			}
			ports = append(ports, convertPort)
		}

		for _, port := range ports {
			s.ShadowsocksRProxy(s.host,
				port,
				nodeInfo.Method,
				nodeInfo.Passwd,
				nodeInfo.Protocol,
				nodeInfo.ProtocolParam,
				nodeInfo.Obfs,
				nodeInfo.ObfsParam,
				nodeInfo.Single,
				&proxy.ShadowsocksRArgs{})
			err := s.Shadowsocksrs[port].Start()
			if err != nil {
				// TODO 错误处理
				panic(err)
			}
		}
	}
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

func (s *shadowsocksrService) AddUser(user *model.UserInfo) error {
	if user2 := s.userTable[user.Uid]; user2 != nil {
		return errors.New(fmt.Sprintf("user %v already exist", user2.Uid))
	}

	if s.nodeInfo.Single == 1 {
		for _, server := range s.Shadowsocksrs {
			server.AddUser(user.Port, user.Passwd)
		}
	} else {
		if s.Shadowsocksrs[user.Port] != nil {
			return errors.New(fmt.Sprintf("add user port %v is used by %v", user.Port, s.PortToUid(user.Port)))
		}
		server := s.ShadowsocksRProxy(s.host,
			user.Port,
			s.nodeInfo.Method,
			user.Passwd,
			s.nodeInfo.Protocol,
			s.nodeInfo.ProtocolParam,
			s.nodeInfo.Obfs,
			s.nodeInfo.ObfsParam,
			s.nodeInfo.Single,
			&proxy.ShadowsocksRArgs{})
		if err := server.Start(); err != nil {
			return errors.Wrap(err, "add user error")
		}
	}
	s.userTable[user.Uid] = user
	return nil
}

func (s *shadowsocksrService) EditUser(user *model.UserInfo) error {
	// TODO after change user profile it will be simultaneously exist old port and new port
	user2 := s.userTable[user.Uid]
	if user2 == nil {
		return errors.New(fmt.Sprintf("user %v dosen't exist", user.Uid))
	}
	if s.nodeInfo.Single != 1 && user.Port != user2.Port && s.Shadowsocksrs[user.Port] != nil{
		return errors.New(fmt.Sprintf("port %v used by user %v",user.Port,s.PortToUid(user.Port)))
	}
	if err := s.DelUser(user.Uid); err != nil {
		return errors.Wrap(err, "edit user del user error")
	}
	if err := s.AddUser(user); err != nil {
		return errors.Wrap(err, "edit user add user error")
	}
	return nil
}

func (s *shadowsocksrService) DelUser(uid int) error {
	port := s.UIDToPort(uid)
	if port == 0 {
		return errors.New(fmt.Sprintf("uid %v is not esixt", uid))
	}

	if s.nodeInfo.Single == 1 {
		for _, server := range s.Shadowsocksrs {
			server.DelUser(port)
			logrus.Infof("server %v del %v success", server.Port, port)
		}
		delete(s.userTable, uid)
		return nil
	} else {
		server := s.Shadowsocksrs[port]
		if server == nil {
			logrus.WithFields(logrus.Fields{
				"port": port,
			}).Info("port is not exist")
			return nil
		}

		if err := server.Close(); err != nil {
			return err
		}
		delete(s.Shadowsocksrs, port)
		delete(s.userTable, uid)
	}
	return nil
}

func (s *shadowsocksrService) GetUserFromPort(port int) *model.UserInfo {
	for _, value := range s.userTable {
		if value.Port == port {
			return value
		}
	}
	return nil
}


func (s *shadowsocksrService) GetUserList() []*model.UserInfo{
	users := make([]*model.UserInfo,0,len(s.userTable))
	for _,value := range s.userTable {
		users = append(users,value)
	}
	return users
}