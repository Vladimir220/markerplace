package db

type IMigrations interface {
	CheckMigrations()
}

type Migrations struct {
}
