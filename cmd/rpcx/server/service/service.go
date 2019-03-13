package service

import (
	"context"
	"time"

	"github.com/rc452860/vnet/cmd/rpcx/model"
	"github.com/rc452860/vnet/common/log"

	"github.com/rc452860/vnet/cmd/rpcx"

	"github.com/pkg/errors"
)

type UserService struct{}

func (userService *UserService) PullEnableUsers(ctx context.Context, arg *rpcx.PullEnableUsersRequest) (*rpcx.PullEnableUsersResponse, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		return nil, errors.WithStack(errors.New("节点验证不通过"))
	}

	users := []model.User{}
	db.Select("id,port,passwd,speed_limit_per_user").
		Where("updated_at >= ?", time.Unix(arg.GetUpdateTime(), 0)).
		Find(&users)
	response := &rpcx.PullEnableUsersResponse{}
	response.EnableUsers = []*rpcx.EnableUser{}
	for _, item := range users {
		tmp := &rpcx.EnableUser{
			Id:       int64(item.ID),
			Port:     uint32(item.Port),
			Password: item.Passwd,
			Limit:    uint64(item.SpeedLimitPerUser),
		}
		response.EnableUsers = append(response.EnableUsers, tmp)
	}
	return response, nil
}

func (userService *UserService) PushUserTraffic(ctx context.Context, arg *rpcx.PushUserTrafficRequest) (*rpcx.PushUserTrafficResponse, error) {
	db, err := GetDB()
	if err != nil {
		log.Err(err)
		return nil, err
	}
	defer db.Close()

	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		return nil, errors.WithStack(errors.New("节点验证不通过"))
	}
	panic("no")
}

func (userService *UserService) PullSsNodeConfig(ctx context.Context, arg *rpcx.SsNodeReuqest) (*rpcx.SsNodeResponse, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ssNode := GetSsNodeServiceInstance().GetSsNodeByIdAndToken(arg.NodeId, arg.Token)
	if ssNode == nil {
		return nil, errors.WithStack(errors.New("节点验证不通过"))
	}

	response := &rpcx.SsNodeResponse{}
	response.Method = ssNode.Method
	return response, nil
}
