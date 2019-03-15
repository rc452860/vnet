package service

type SsNodeOnlineService struct {
}

var ssNodeOnlineServiceInstance = NewSsNodeOnlineService()

func NewSsNodeOnlineService() *SsNodeOnlineService {
	return &SsNodeOnlineService{}
}

func GetSsNodeOnlineServiceInstance() *SsNodeOnlineService {
	return ssNodeOnlineServiceInstance
}

func (s *SsNodeOnlineService) RecordSsNodeOnline(nodeId int64, onlineUsers int32) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	db.Exec("INSERT INTO `ss_node_online_log` (id,node_id,online_user,log_time) "+
		" VALUES (NULL,?,?,unix_timestamp())", nodeId, onlineUsers)
	return nil
}
