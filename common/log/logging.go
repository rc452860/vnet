package log

import "github.com/sirupsen/logrus"

// default use root logging
func Debug(message string, params ...interface{}) {
	logrus.Debugf(message,params)
}

// default use root logging
func Info(message string, params ...interface{}) {
	logrus.Infof(message,params)
}

// default use root logging
func Warn(message string, params ...interface{}) {
	logrus.Warnf(message,params)
}

// default use root logging
func Error(message string, params ...interface{}) {
	logrus.Errorf(message,params)
}

// default use root logging
func Err(err error) {
	logrus.Errorf("%v",err)
}