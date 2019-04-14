package skr

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func GetLogger() *logrus.Logger {
	return logger
}
