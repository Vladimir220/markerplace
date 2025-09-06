package logger_lib

import (
	"context"
	"fmt"
)

func CreateLoggerGateway(ctx context.Context, parentName string) ILogger {
	remoteLogger, err := CreateKafkaLogger(ctx, parentName)
	localLogger := CreateLocalLogger(parentName)
	var remoteUnavailable bool
	if err != nil {
		localLogger.WriteWarning(fmt.Sprintf("%s: %v", "CreateLoggerGateway(): remote logger unavailable", err))
		remoteUnavailable = true
	}

	config := GetLoggerConfig()

	return &LoggerGateway{
		remoteUnavailable: remoteUnavailable,
		localLogger:       localLogger,
		remoteLogger:      remoteLogger,
		ctx:               ctx,
		config:            config,
	}
}

// Proxy for ILogger
// Adapter for ILoggerRemote
type LoggerGateway struct {
	remoteUnavailable bool
	config            LoggerConfig
	localLogger       ILogger
	remoteLogger      ILoggerRemote
	ctx               context.Context
}

func (l LoggerGateway) WriteWarning(msg string) {
	if l.config.WarningsToStdOut {
		fmt.Println(msg)
	}

	var err error
	if !l.remoteUnavailable {
		err = l.remoteLogger.WriteWarning(msg)
	}
	if l.remoteUnavailable || err != nil {
		l.localLogger.WriteWarning(fmt.Sprintf("%s: %v", "WriteWarning(): remote logger unavailable", err))
		l.localLogger.WriteWarning(msg)
	}
}

func (l LoggerGateway) WriteError(msg string) {
	if l.config.ErrorsToStdOut {
		fmt.Println(msg)
	}

	var err error
	if !l.remoteUnavailable {
		err = l.remoteLogger.WriteError(msg)
	}
	if l.remoteUnavailable || err != nil {
		l.localLogger.WriteWarning(fmt.Sprintf("%s: %v", "WriteError(): remote logger unavailable", err))
		l.localLogger.WriteError(msg)
	}
}

func (l LoggerGateway) WriteInfo(msg string) {
	if l.config.InfoToStdOut {
		fmt.Println(msg)
	}

	var err error
	if !l.remoteUnavailable {
		err = l.remoteLogger.WriteInfo(msg)
	}
	if l.remoteUnavailable || err != nil {
		l.localLogger.WriteWarning(fmt.Sprintf("%s: %v", "WriteInfo(): remote logger unavailable", err))
		l.localLogger.WriteInfo(msg)
	}
}
