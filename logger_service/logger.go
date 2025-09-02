package main

import (
	"fmt"
	"log"
	"os"
)

var ErrorLogPath = "./logs/error.log"
var InfoLogPath = "./logs/info.log"
var WarningLogPath = "./logs/warning.log"

type ILogger interface {
	WriteWarning(msg string)
	WriteError(msg string)
	WriteInfo(msg string)
}

func CreateLogger(parentName string) ILogger {
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

	logger := &Logger{
		warning: log.New(warningLogFile, fmt.Sprintf("WARNING:%s:", parentName), log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		err:     log.New(errorLogFile, fmt.Sprintf("ERROR:%s:", parentName), log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
		info:    log.New(infoLogFile, fmt.Sprintf("INFO:%s:", parentName), log.Ldate|log.Ltime|log.Lshortfile|log.Lmsgprefix),
	}

	return logger
}

type Logger struct {
	warning *log.Logger
	err     *log.Logger
	info    *log.Logger
}

func (l Logger) WriteWarning(msg string) {
	l.warning.Println(msg)
}

func (l Logger) WriteError(msg string) {
	l.err.Println(msg)
}

func (l Logger) WriteInfo(msg string) {
	l.info.Println(msg)
}
