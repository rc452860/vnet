package model

type NodeInfo struct {
	ID            int    `json:"id"`
	Method        string `json:"method"`
	Protocol      string `json:"protocol"`
	Obfs          string `json:"obfs"`
	ProtocolParam string `json:"protocol_param"`
	ObfsParam     string `json:"obfs_param"`
}

type UserInfo struct {
	Uid    int    `json:"uid"`
	Port   int    `json:"port"`
	Passwd string `json:"passwd"`
	Limit  uint64 `json:"speed_limit_per_user"`
	Enable int    `json:"enable"`
}

type UserTraffic struct {
	Uid      int   `json:"uid"`
	Upload   int64 `json:"upload"'`
	Download int64 `json:"download"`
}

type NodeOnline struct {
	Uid int    `json:"uid"`
	IP  string `json:"ip"`
}

type NodeStatus struct {
	CPU string `json:"cpu"`
	MEM string `json:"mem"`
	NET string `json:"net"`
}
