package logger_lib

import (
	"fmt"
	"log"
	"os"
)

var ErrorLogPath = "./logs/error.log"
var InfoLogPath = "./logs/info.log"
var WarningLogPath = "./logs/warning.log"

func CreateLocalLogger(parentName string) ILogger {
	errorLogFile, err := os.OpenFile(ErrorLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	infoLogFile, err := os.OpenFile(InfoLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	warningLogFile, err := os.OpenFile(WarningLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}

	LocalLogger := &LocalLogger{
		warning: log.New(warningLogFile, fmt.Sprintf("WARNING:%s:", parentName), log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		err:     log.New(errorLogFile, fmt.Sprintf("ERROR:%s:", parentName), log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		info:    log.New(infoLogFile, fmt.Sprintf("INFO:%s:", parentName), log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
	}

	return LocalLogger
}

type LocalLogger struct {
	warning *log.Logger
	err     *log.Logger
	info    *log.Logger
}

func (l LocalLogger) WriteWarning(msg string) {
	l.warning.Println(msg)
}

func (l LocalLogger) WriteError(msg string) {
	l.err.Println(msg)
}

func (l LocalLogger) WriteInfo(msg string) {
	l.info.Println(msg)
}
