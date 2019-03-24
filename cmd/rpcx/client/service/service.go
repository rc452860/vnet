package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/rc452860/vnet/record"
	"github.com/rc452860/vnet/service"

	"github.com/rc452860/vnet/component/cachex"

	"github.com/rc452860/vnet/utils/goroutine"

	"github.com/rc452860/vnet/cmd/rpcx"
	"github.com/rc452860/vnet/cmd/rpcx/config"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/proxy/server"
	. "github.com/rc452860/vnet/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var ssNode *rpcx.SsNodeResponse
var timeout = 3 * time.Second

func Start() {
	// 获取节点配置
	var err error
	ssNode, err = PullSsNodeConfig()
	if err != nil {
		log.Error("get ssNode config error: %v", err)
		return
	}

	log.Info("get node config:( %v )", ssNode)

	// fetch user in time wheel
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go goroutine.Protect(func() {
		defer wg.Done()
		VnetStart()
	})

	go goroutine.Protect(func() {
		HttpStart()
	})
	wg.Wait()
}

var (
	g_updatetime int64 = 0
)

// start all task
func VnetStart() {
	tick := time.Tick(1 * time.Second)
	count := 0
	record.GetGRMInstance()
	var syncTime, addtionSyncTime int
	syncTime = StringToInt(viper.GetString(config.C_SyncInterval))
	addtionSyncTime = StringToInt(viper.GetString(config.C_ReportInterval))
	for {
		<-tick
		count++
		// database synchronize time wheel
		if count%syncTime == 0 {
			err := VnetTask()
			if err != nil {
				log.Err(err)
			}
		}

		// one minuts time wheel
		if count%addtionSyncTime == 0 {
			err := VnetTrafficTask()
			if err != nil {
				log.Err(err)
			}
			err = VnetOnlineUserNumberTask()
			if err != nil {
				log.Err(err)
			}
			err = VnetSsNodeInfoTask()
			if err != nil {
				log.Err(err)
			}
		}
	}
}

// start shadowsocks service
func VnetTask() error {
	log.Info("emit enbale user synchronize task time: %v", g_updatetime)
	r, err := GetUsers()
	if err != nil {
		return err
	}
	if len(r.EnableUsers) > 1 {
		log.Info("recive users len:%v", len(r.EnableUsers))
	} else {
		log.Debug("recive users len:%v", len(r.EnableUsers))
	}
	for _, item := range r.EnableUsers {
		ssservice := CurrentShadowsocksService()
		// start new service if port not exist,and restart config changed port
		if item.Enable {
			ssservice.AddAndStart("0.0.0.0", ssNode.Method, item.Password, int(item.Port), server.ShadowsocksArgs{
				ConnectTimeout: 0,
				Limit:          GetFinalLimit(item.Limit) * 1024,
				TCPSwitch:      viper.GetString(config.C_TCP),
				UDPSwitch:      viper.GetString(config.C_UDP),
				Data:           item,
			})
		}
		// stop user
		if !item.Enable {
			ssservice.Del(int(item.Port))
		}
	}

	return nil
}

// VnetTrafficTask upload vnet traffic
func VnetTrafficTask() error {
	log.Info("emit traffic report task")
	ssService := CurrentShadowsocksService()
	cacheService := cachex.GetCache()

	vnetList := ssService.List()
	trafficInfos := []*rpcx.TrafficInfo{}

	for _, item := range vnetList {
		user, ok := item.ShadowsocksArgs.Data.(*rpcx.EnableUser)
		if !ok {
			log.Error("data convert to user failed at traffic task on %v", item.Port)
			continue
		}
		upKey := fmt.Sprintf("traffic_task,up,%v", item.Port)
		downKey := fmt.Sprintf("traffic_task,down,%v", item.Port)
		var lastUp, lastDown, currentUp, currentDown uint64
		if r, ok := cacheService.Get(upKey).(uint64); ok {
			lastUp = r
		}
		if r, ok := cacheService.Get(downKey).(uint64); ok {
			lastDown = r
		}

		currentUp = item.UpBytes
		currentDown = item.DownBytes
		cacheService.Put(upKey, currentUp, 10*time.Minute)
		cacheService.Put(downKey, currentDown, 10*time.Minute)
		// filter 500k
		if currentUp+currentDown < lastUp+lastDown+500*1024 {
			continue
		}
		tmp := &rpcx.TrafficInfo{
			Uid:      user.Id,
			Port:     int64(item.Port),
			Upload:   currentUp - lastUp,
			Download: currentDown - lastDown,
		}

		trafficInfos = append(trafficInfos, tmp)
	}
	if len(trafficInfos) > 0 {
		log.Info("traiifc task report data length:%v", len(trafficInfos))
	} else {
		log.Debug("traiifc task report data length:%v", len(trafficInfos))
	}
	_, err := PushUserTraffic(trafficInfos)
	return err
}

func VnetOnlineUserNumberTask() error {
	log.Info("emit online users count report task")
	grmInstance := record.GetGRMInstance()
	_, err := PushSsNodeOnlineLog(grmInstance.GetLastOneMinuteOnlineCount())
	log.Info("current users count is:%v", grmInstance.GetLastOneMinuteOnlineCount())
	return err
}

func VnetSsNodeInfoTask() error {
	log.Info("emit node info report task")
	_, err := PushSsNodeInfo(formatLoad())
	return err
}

// ConOpStr condition operation
func ConOpStr(cond bool, a string, b string) string {
	if cond {
		return a
	}
	return b
}

// StringToInt convert str to int
func StringToInt(data string) int {
	r, e := strconv.Atoi(data)
	if e != nil {
		return 0
	}
	return r
}

// GetUsers is wrap rpcx
func GetUsers() (*rpcx.PullEnableUsersResponse, error) {
	conn, err := grpc.Dial(viper.GetString(config.C_RpcAddress), grpc.WithTimeout(1*time.Second), grpc.WithInsecure(), grpc.WithTimeout(3*time.Second), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := rpcx.NewUserServiceClient(conn)
	defer func() {
		g_updatetime = time.Now().Unix()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.PullEnableUsers(ctx, &rpcx.PullEnableUsersRequest{
		NodeId:     int64(StringToInt(viper.GetString(config.C_NodeId))),
		Token:      viper.GetString(config.C_Token),
		UpdateTime: g_updatetime,
	})
}

//PullSsNodeConfig is wrap rpcx
func PullSsNodeConfig() (*rpcx.SsNodeResponse, error) {
	conn, err := grpc.Dial(viper.GetString(config.C_RpcAddress), grpc.WithTimeout(1*time.Second), grpc.WithInsecure(), grpc.WithTimeout(3*time.Second), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := rpcx.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.PullSsNodeConfig(ctx, &rpcx.SsNodeReuqest{
		NodeId: int64(StringToInt(viper.GetString(config.C_NodeId))),
		Token:  viper.GetString(config.C_Token),
	})
}

// PushUserTraffic is wrap rpcx
func PushUserTraffic(data []*rpcx.TrafficInfo) (*rpcx.PushUserTrafficResponse, error) {
	conn, err := grpc.Dial(viper.GetString(config.C_RpcAddress), grpc.WithTimeout(1*time.Second), grpc.WithInsecure(), grpc.WithTimeout(3*time.Second), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := rpcx.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := &rpcx.PushUserTrafficRequest{
		NodeId:      int64(StringToInt(viper.GetString(config.C_NodeId))),
		Token:       viper.GetString(config.C_Token),
		TrafficInfo: data,
	}

	return c.PushUserTraffic(ctx, request)
}

func PushSsNodeOnlineLog(onlineUserNum int) (*rpcx.SsNodeOnlineLogResponse, error) {
	conn, err := grpc.Dial(viper.GetString(config.C_RpcAddress), grpc.WithTimeout(1*time.Second), grpc.WithInsecure(), grpc.WithTimeout(3*time.Second), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := rpcx.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := &rpcx.SsNodeOnlineLogRequest{
		NodeId:     int64(StringToInt(viper.GetString(config.C_NodeId))),
		Token:      viper.GetString(config.C_Token),
		OnlineUser: int32(onlineUserNum),
	}
	return c.PushSsNodeOnlineLog(ctx, request)

}

func PushSsNodeInfo(load string) (*rpcx.SsNodeInfoResponse, error) {
	conn, err := grpc.Dial(viper.GetString(config.C_RpcAddress), grpc.WithTimeout(1*time.Second), grpc.WithInsecure(), grpc.WithTimeout(3*time.Second), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := rpcx.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := &rpcx.SsNodeInfoRequest{
		NodeId: int64(StringToInt(viper.GetString(config.C_NodeId))),
		Token:  viper.GetString(config.C_Token),
		Load:   load,
	}
	return c.PushSsNodeInfo(ctx, request)
}

func formatLoad() string {
	return fmt.Sprintf("cpu:%v%% mem:%v%% disk:%v%%", service.GetCPUUsage(), service.GetMemUsage(), service.GetDiskUsage())
}

func GetFinalLimit(limit uint64) uint64 {
	result, err := strconv.ParseInt(viper.GetString(config.C_LIMIT), 10, 64)
	if err != nil {
		log.Error("config limit error:%v", err)
		return limit
	}
	configLimit := uint64(result)
	if configLimit == 0 || limit < configLimit {
		return limit
	}
	return configLimit
}
