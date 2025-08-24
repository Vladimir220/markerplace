package postgres

import (
	"database/sql"
	"fmt"
	"main/models"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type IMarketplaceDAO interface {
	GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error)
	GetUser(login string) (user models.User, password string, isFound bool, err error)
	NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error)
	Registr(login, password string) (user models.User, isAlreadyExist bool, err error)
	Close()
}

func CreateMarketplaceDAO() (dao IMarketplaceDAO, err error) {
	logLabel := "CreateMarketplaceDAO():"

	err = checkDbExistence()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	connection, err := connect()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	mpDao := &MarketplaceDAO{
		connection: connection,
	}

	err = mpDao.checkMigrations()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	return mpDao, nil
}

type MarketplaceDAO struct {
	connection *sql.DB
}

func (md MarketplaceDAO) Close() {
	md.connection.Close()
}

func (md *MarketplaceDAO) checkMigrations() error {
	logLabel := "checkMigrations():"

	if md.connection == nil {
		return fmt.Errorf("%s no connection", logLabel)
	}

	driver, err := postgres.WithInstance(md.connection, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s checkMigrations(): %v", logLabel, err)
	}

	dbName := os.Getenv("DB_NAME")
	var path string
	currentDir, _ := os.Getwd()
	path = filepath.ToSlash(currentDir)
	if path != "" {
		path = path + "/"
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path+"db/migrations", dbName, driver)
	if err != nil {
		return fmt.Errorf("%s checkMigrations(): %v", logLabel, err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s checkMigrations(): %v", logLabel, err)
	}

	return nil
}
