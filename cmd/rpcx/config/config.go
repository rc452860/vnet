package config

const (
	C_NodeId         = "nodeId"
	C_Token          = "token"
	C_RpcAddress     = "rpcAddress"
	C_SyncInterval   = "syncInterval"
	C_ReportInterval = "addtionSyncTime"
	C_TCP            = "tcp"
	C_UDP            = "udp"
)

type FlagSetting struct {
	Name     string
	Default  string
	Usage    string
	Example  string
	Required bool
}

var ClientConfig = []FlagSetting{
	FlagSetting{
		Name:     C_NodeId,
		Usage:    "node",
		Required: true,
	},
	FlagSetting{
		Name:     C_Token,
		Usage:    "token",
		Required: true,
	},
	FlagSetting{
		Name:     C_RpcAddress,
		Usage:    "rpcaddress example:0.0.0.0:5050",
		Required: true,
	},
	FlagSetting{
		Name:     C_SyncInterval,
		Usage:    "sync time unit(second) default 2",
		Default:  "2",
		Required: true,
	},
	FlagSetting{
		Name:     C_ReportInterval,
		Usage:    "user traffic, server info, online user sync time",
		Default:  "60",
		Required: true,
	},
	FlagSetting{
		Name:     C_TCP,
		Usage:    "allow tcp connection",
		Default:  "true",
		Required: true,
	},
	FlagSetting{
		Name:     C_UDP,
		Usage:    "allow UDP connection",
		Default:  "true",
		Required: true,
	},
}

const (
	// database source
	S_DS      = "ds"
	S_RPCPort = "rpcPort"
)

var ServerConfig = []FlagSetting{
	FlagSetting{
		Name:     S_DS,
		Usage:    "database source ([username]:[password]@tcp([ip or url]:[port])/[database]?parseTime=true&loc=UTC) example: root:root@tcp(localhost:3306)/ssrpanel?parseTime=true&loc=UTC",
		Required: true,
	},
	FlagSetting{
		Name:     S_RPCPort,
		Usage:    "rpc listen port",
		Required: true,
	},
}
