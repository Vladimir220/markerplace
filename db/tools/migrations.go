package tools

type IMigrations interface {
	CheckMigrations()
}

type Migrations struct {
}
