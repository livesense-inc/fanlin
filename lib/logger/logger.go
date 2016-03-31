package logger

import (
	"os"

	"github.com/mymch/formatter"
	"github.com/sirupsen/logrus"
)

func NewLogger(path string) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(formatter.SysLogFormatter)
	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalln("Can not create log file: ", path)
	}
	logger.Out = logFile
	return logger
}
