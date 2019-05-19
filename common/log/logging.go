package log

import "github.com/sirupsen/logrus"

// default use root logging
func Debug(message string, params ...interface{}) {
	if params != nil {
		logrus.Debugf(message, params)
	} else {
		logrus.Debug(message)
	}
}

// default use root logging
func Info(message string, params ...interface{}) {
	if params != nil {
		logrus.Infof(message, params)
	} else {
		logrus.Info(message)
	}
}

// default use root logging
func Warn(message string, params ...interface{}) {
	if params != nil {
		logrus.Warnf(message, params)
	} else {
		logrus.Warn(message)
	}
}

// default use root logging
func Error(message string, params ...interface{}) {
	if params != nil {
		logrus.Error(message, params)
	} else {
		logrus.Errorf(message)
	}
}

// default use root logging
func Err(err error) {
	logrus.Errorf("%v", err)
}
