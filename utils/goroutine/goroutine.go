package goroutine

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func Protect(g func()) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("run time panic: %s stack: %s", err, string(debug.Stack()))
		}
	}()
	g()
}
