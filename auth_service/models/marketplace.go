package models

type User struct {
	Login string `redis:"login"`
	Group string `redis:"group"`
}

type LogsConfig struct {
	InfoLogs              bool
	PrintErrorsToStdOut   bool
	PrintWarningsToStdOut bool
	PrintInfoToStdOut     bool
}
