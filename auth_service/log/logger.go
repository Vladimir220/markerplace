package log

type ILogger interface {
	WriteWarning(msg string)
	WriteError(msg string)
	WriteInfo(msg string)
}

type ILoggerRemote interface {
	WriteWarning(msg string) (err error)
	WriteError(msg string) (err error)
	WriteInfo(msg string) (err error)
}
