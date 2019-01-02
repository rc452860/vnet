package db

import (
	"testing"

	"github.com/rc452860/vnet/log"
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
