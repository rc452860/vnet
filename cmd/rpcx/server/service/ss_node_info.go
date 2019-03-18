package service

type SsNodeInfo struct {
}

var ssNodeInfo = NewSsNodeInfo()

func NewSsNodeInfo() *SsNodeInfo {
	return &SsNodeInfo{}
}

func GetSsNodeInfoInstance() *SsNodeInfo {
	return ssNodeInfo
}

func (s *SsNodeInfo) RecordSsNodeInfo(nodeId int64, load string) error {

	db, err := GetDB()
	if err != nil {
		return err
	}

	db.Exec("INSERT INTO `ss_node_info` (`id`,`node_id`,`uptime`,`load`,`log_time`)"+
		"VALUES (NULL,?,unix_timestamp(),?,unix_timestamp())", nodeId, load)
	return nil
}
