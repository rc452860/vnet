package conn

import "github.com/rc452860/vnet/log"

var logging *log.Logging

func init() {
	logging = log.GetLogger("root")
}
