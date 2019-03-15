package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rc452860/vnet/component/cachex"

	"github.com/rc452860/vnet/cmd/rpcx/model"
	"github.com/rc452860/vnet/common/log"

	"github.com/rc452860/vnet/cmd/rpcx"

	"github.com/pkg/errors"
)

type UserService struct{}

func (userService *UserService) PullEnableUsers(ctx context.Context, arg *rpcx.PullEnableUsersRequest) (*rpcx.PullEnableUsersResponse, error) {
	log.Info("%v,%s", arg.NodeId, arg.Token)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		log.Error("node %v token verification failed,error token: %s", arg.NodeId, arg.Token)
		return nil, errors.WithStack(errors.New("node verification failed"))
	}

	users := []model.User{}

	date := time.Now()
	if userService.getUpdateTime(arg.NodeId, arg.Token).Unix() != time.Unix(0, 0).Unix() && arg.UpdateTime != 0 {
		fmt.Println(userService.getUpdateTime(arg.NodeId, arg.Token))
		db.Select("id,port,passwd,speed_limit_per_user,enable").
			Where("updated_at >= ?", userService.getUpdateTime(arg.NodeId, arg.Token).Add(-3*time.Second)).
			Find(&users)
	} else {
		db.Select("id,port,passwd,speed_limit_per_user,enable").
			Where("enable = 1").
			Find(&users)
	}
	userService.setUpdateTime(arg.NodeId, arg.Token, date)

	response := &rpcx.PullEnableUsersResponse{}
	response.EnableUsers = []*rpcx.EnableUser{}

	for _, item := range users {
		tmp := &rpcx.EnableUser{
			Id:       int64(item.ID),
			Port:     uint32(item.Port),
			Password: item.Passwd,
			Limit:    uint64(item.SpeedLimitPerUser),
			Enable:   item.Enable == 1,
		}
		response.EnableUsers = append(response.EnableUsers, tmp)
	}

	log.Info("read len:%v", len(response.EnableUsers))
	return response, nil
}

func (userService *UserService) PushUserTraffic(ctx context.Context, arg *rpcx.PushUserTrafficRequest) (*rpcx.PushUserTrafficResponse, error) {
	log.Info("%v,%s", arg.NodeId, arg.Token)
	db, err := GetDB()
	if err != nil {
		log.Err(err)
		return nil, err
	}
	defer db.Close()

	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		log.Error("node %v token verification failed,error token: %s", arg.NodeId, arg.Token)
		return nil, errors.WithStack(errors.New("node verification failed"))
	}
	log.Info("nodeId: %v,len: %v", arg.NodeId, len(arg.TrafficInfo))
	GetTrafficServiceInstance().RecordTrafficLogSync(arg)
	return &rpcx.PushUserTrafficResponse{}, nil
}

func (userService *UserService) PullSsNodeConfig(ctx context.Context, arg *rpcx.SsNodeReuqest) (*rpcx.SsNodeResponse, error) {
	log.Info("%v,%s", arg.NodeId, arg.Token)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		log.Error("node %v token verification failed,error token: %s", arg.NodeId, arg.Token)
		return nil, errors.WithStack(errors.New("node verification failed"))
	}

	response := &rpcx.SsNodeResponse{
		Method:      ssNode.Method,
		Tcp:         int32(ssNode.Tcp),
		Udp:         int32(ssNode.Udp),
		TrafficRate: int32(ssNode.TrafficRate),
	}
	log.Info("%v", response)
	return response, nil
}

func (userService *UserService) PushSsNodeOnlineLog(ctx context.Context, arg *rpcx.SsNodeOnlineLogRequest) (*rpcx.SsNodeOnlineLogResponse, error) {
	log.Info("%v,%s", arg.NodeId, arg.Token)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		log.Error("node %v token verification failed,error token: %s", arg.NodeId, arg.Token)
		return nil, errors.WithStack(errors.New("node verification failed"))
	}

	s := GetSsNodeOnlineServiceInstance()
	err = s.RecordSsNodeOnline(arg.NodeId, arg.OnlineUser)
	if err != nil {
		log.Err(err)
	}
	return &rpcx.SsNodeOnlineLogResponse{}, nil
}

func (userService *UserService) PushSsNodeInfo(ctx context.Context, arg *rpcx.SsNodeInfoRequest) (*rpcx.SsNodeInfoResponse, error) {
	log.Info("%v,%s", arg.NodeId, arg.Token)
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		log.Error("node %v token verification failed,error token: %s", arg.NodeId, arg.Token)
		return nil, errors.WithStack(errors.New("node verification failed"))
	}

	s := GetSsNodeInfoInstance()
	err = s.RecordSsNodeInfo(arg.NodeId, arg.Load)
	if err != nil {
		log.Err(err)
	}
	return &rpcx.SsNodeInfoResponse{}, nil
}

/*--------------------------------------cute split line---------------------------------------------*/
func (userService *UserService) getUpdateTime(nodeid int64, token string) time.Time {
	key := fmt.Sprintf("updatetime,%v,%s", nodeid, token)
	result := cachex.GetCache().Get(key)
	if r, ok := result.(time.Time); ok {
		return r
	}
	return time.Unix(0, 0)
}

func (userService *UserService) setUpdateTime(nodeid int64, token string, date time.Time) {
	key := fmt.Sprintf("updatetime,%v,%s", nodeid, token)
	cachex.GetCache().Put(key, date, 1*time.Minute)
}
