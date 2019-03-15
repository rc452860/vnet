package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rc452860/vnet/cmd/rpcx/model"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/component/cachex"
)

var ssNodeServiceInstance = NewSsNodeService()

type SsNodeService struct {
}

func NewSsNodeService() *SsNodeService {
	return &SsNodeService{}
}

// NewSsNodeService 单例模式
func GetSsNodeServiceInstance() *SsNodeService {
	return ssNodeServiceInstance
}

// 验证节点
func (ssNodeService *SsNodeService) VerifyNode(nodeId int64, token string) bool {
	if ssNodeService.GetSsNodeByIdAndToken(nodeId, token) != nil {
		return true
	}
	return false
}

// 通过id和token获得节点配置信息
func (ssNodeService *SsNodeService) GetSsNodeByIdAndToken(nodeId int64, token string) *model.SsNode {
	cacheKey := fmt.Sprintf("%s,%s", strconv.Itoa(int(nodeId)), token)
	cacheExpire := 5 * time.Minute
	if result, ok := cachex.GetCache().Get(cacheKey).(*model.SsNode); ok {
		return result
	}
	db, err := GetDB()
	if err != nil {
		log.Err(err)
		return nil
	}
	ssNode := &model.SsNode{}
	db = db.Select("id,token,method,traffic_rate,tcp,udp").Where("id=? and token=?", int(nodeId), token).First(ssNode)
	if db.RecordNotFound() {
		return nil
	}
	if db.GetErrors() != nil && len(db.GetErrors()) > 0 {
		for _, item := range db.GetErrors() {
			if item.Error() != "record not found" {
				log.Err(item)
			}
		}
		return nil
	}
	cachex.GetCache().Put(cacheKey, ssNode, cacheExpire)
	return ssNode
}
