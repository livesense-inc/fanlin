package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/ieee0824/logrus-formatter"
)

func NewLogger(path string) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = new(formatter.SysLogFormatter)
	logFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalln("Can not create log file: ", path)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logger.Out = mw
	return logger
}
