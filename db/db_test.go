package db

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rc452860/vnet/comm/log"
)

func Test_GetEnableUser(t *testing.T) {
	logging := log.GetLogger("root")

	ssNodeList, err := GetEnableUser()
	if err != nil {
		t.Error(err)
	}
	for _, item := range ssNodeList {
		logging.Info("item: %v", item)
	}
}

func TestFormatLoad(t *testing.T) {
	t.Log(FormatLoad())
}
