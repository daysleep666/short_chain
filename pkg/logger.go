package pkg

import (
	"fmt"
	"os"

	"github.com/labstack/gommon/log"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var LOGGER_INSTANCE Logger

func InitLogger(filePath string) (err error) {
	f, err := os.Create(filePath)
	if err != nil {
		return
	}
	log.SetOutput(f)
	log.SetLevel(log.DEBUG)
	log.EnableColor()
	return
}

func NewLogger(loggid string) Logger {
	return log.New(fmt.Sprintf("[logid:%s]", loggid))
}
