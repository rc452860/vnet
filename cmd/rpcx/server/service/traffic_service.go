package service

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/rc452860/vnet/cmd/rpcx"

	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/utils/datasize"
)

var trafficServiceInstance = NewTrafficService()

type TrafficService struct {
	sync.Mutex
	TrafficChan chan *rpcx.PushUserTrafficRequest
	cancel      context.CancelFunc
	ctx         context.Context
	started     bool
}

func NewTrafficService() *TrafficService {
	instance := &TrafficService{
		TrafficChan: make(chan *rpcx.PushUserTrafficRequest, 1024),
	}
	instance.Start()
	return instance
}

func GetTrafficServiceInstance() *TrafficService {
	return trafficServiceInstance
}

func (trafficService *TrafficService) Start() {
	trafficService.Lock()
	defer trafficService.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	trafficService.ctx = ctx
	trafficService.cancel = cancel
	trafficService.started = true
	go trafficService.Monitor()
}

func (trafficService *TrafficService) Stop() {
	trafficService.Lock()
	defer trafficService.Unlock()
	trafficService.cancel()
	<-trafficService.ctx.Done()
}

// TODO destroy

// Monitor 监视者
func (trafficService *TrafficService) Monitor() {
	for {
		var data *rpcx.PushUserTrafficRequest
		select {
		case <-trafficService.ctx.Done():
			return
		case data = <-trafficService.TrafficChan:
		}
		trafficService.RecordTrafficLog(data)
	}
}

// RecordTrafficLogSync 异步记录日志
func (trafficService *TrafficService) RecordTrafficLogSync(data *rpcx.PushUserTrafficRequest) {
	trafficService.TrafficChan <- data
}

// RecordTrafficLog 同步记录日志
func (trafficService *TrafficService) RecordTrafficLog(data *rpcx.PushUserTrafficRequest) {
	// 拿到db
	db, err := GetDB()
	if err != nil {
		log.Err(err)
		return
	}
	defer db.Close()
	// 拿到ssNode配置
	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(data.NodeId, data.Token)
	// 批量执行插入操作
	buf := new(bytes.Buffer)
	whenUp := new(bytes.Buffer)
	whenDown := new(bytes.Buffer)
	whenPort := new(bytes.Buffer)
	for _, item := range data.TrafficInfo {
		if item.Upload+item.Download == 0 {
			continue
		}
		traffic := datasize.MustHumanSize(uint64(item.Download+item.Upload) * uint64(ssNode.TrafficRate))
		whenUp.WriteString(fmt.Sprintf(" WHEN %v THEN u+%v", item.Port, float32(item.Download)*ssNode.TrafficRate))
		whenDown.WriteString(fmt.Sprintf(" WHEN %v THEN d+%v", item.Port, float32(item.Upload)*ssNode.TrafficRate))
		whenPort.WriteString(fmt.Sprintf("%v,", item.Port))

		buf.WriteString(fmt.Sprintf("(NULL,%v,%v,%v,%v,%f,'%s',unix_timestamp()),",
			item.Uid,
			item.Download,
			item.Upload,
			data.NodeId,
			ssNode.TrafficRate,
			traffic))
	}
	if buf.Len() <= 1 {
		return
	}
	db.Exec("INSERT INTO `user_traffic_log` (`id`, `user_id`, `u`, `d`, " +
		"`node_id`, `rate`, `traffic`, `log_time`) VALUES" + buf.String()[:buf.Len()-1])

	// 上传和下载相反
	updateUserTrafficSQL := fmt.Sprintf("UPDATE user SET u = CASE port%s END,"+
		"d = CASE port%s END,t = unix_timestamp() WHERE port IN (%s)",
		whenUp.String(),
		whenDown.String(),
		whenPort.String()[:whenPort.Len()-1])
	tx := db.Begin()
	tx.Exec(updateUserTrafficSQL)
	tx.Commit()
	// TODO 如果有报错打印日志 暂时忽略错误
	if db.GetErrors() != nil && len(db.GetErrors()) > 0 {
		for _, item := range db.GetErrors() {
			log.Err(item)
		}
	}
}
