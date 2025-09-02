package log

import (
	"context"
	"fmt"
)

func CreateLoggerAdapter(ctx context.Context, parentName string) ILogger {
	remoteLogger, err := CreateKafkaLogger(ctx, parentName)
	localLogger := CreateLocalLogger(parentName)
	var remoteUnavailable bool
	if err != nil {
		localLogger.WriteWarning(fmt.Sprintf("%s: %v", "CreateLoggerAdapter(): remote logger unavailable", err))
		remoteUnavailable = true
	}

	return &LoggerAdapter{
		remoteUnavailable: remoteUnavailable,
		localLogger:       localLogger,
		remoteLogger:      remoteLogger,
		ctx:               ctx,
	}
}

type LoggerAdapter struct {
	remoteUnavailable bool
	localLogger       ILogger
	remoteLogger      ILoggerRemote
	ctx               context.Context
}

func (l LoggerAdapter) WriteWarning(msg string) {
	var err error
	if !l.remoteUnavailable {
		err = l.remoteLogger.WriteWarning(msg)
	}
	if l.remoteUnavailable || err != nil {
		l.localLogger.WriteWarning(fmt.Sprintf("%s: %v", "WriteWarning(): remote logger unavailable", err))
		l.localLogger.WriteWarning(msg)
	}
}

func (l LoggerAdapter) WriteError(msg string) {
	var err error
	if !l.remoteUnavailable {
		err = l.remoteLogger.WriteError(msg)
	}
	if l.remoteUnavailable || err != nil {
		l.localLogger.WriteWarning(fmt.Sprintf("%s: %v", "WriteError(): remote logger unavailable", err))
		l.localLogger.WriteError(msg)
	}
}

func (l LoggerAdapter) WriteInfo(msg string) {
	var err error
	if !l.remoteUnavailable {
		err = l.remoteLogger.WriteInfo(msg)
	}
	if l.remoteUnavailable || err != nil {
		l.localLogger.WriteWarning(fmt.Sprintf("%s: %v", "WriteInfo(): remote logger unavailable", err))
		l.localLogger.WriteInfo(msg)
	}
}
