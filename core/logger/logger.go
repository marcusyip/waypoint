package logger

import "github.com/Sirupsen/logrus"

var logger = logrus.New()

func GetLogger() *logrus.Logger {
	return logger
}
