package command

const (
	API_HOST = "api_host"
	HOST     = "host"
	NODE_ID  = "node_id"
	KEY      = "key"
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
		Name:     API_HOST,
		Usage:    "api host example: http://localhost",
		Required: true,
	},
	FlagSetting{
		Name:     HOST,
		Usage:    "host example: 0.0.0.0",
		Required: true,
		Default: "0.0.0.0",
	},
	FlagSetting{
		Name:     NODE_ID,
		Usage:    "node_id",
		Required: true,
	},
	FlagSetting{
		Name:     KEY,
		Usage:    "key",
		Required: true,
	},
}
