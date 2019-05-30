package command

const (
	NODE_ID = "nodeId"
	KEY     = "key"
)

type FlagSetting struct {
	Name     string
	Default  string
	Usage    string
	Example  string
	Required bool
}

var flagConfigs = []FlagSetting{
	FlagSetting{
		Name:     NODE_ID,
		Usage:    "nodeid",
		Required: true,
	},
	FlagSetting{
		Name:     KEY,
		Usage:    "key",
		Required: true,
	},
}
