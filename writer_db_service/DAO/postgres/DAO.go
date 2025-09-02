package postgres

import (
	"database/sql"
	"fmt"
	"writer_db_service/models"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type IWriterMarketplaceDAO interface {
	NewAnnouncement(announcement models.ExtendedAnnouncement) (err error)
	Close()
}

func CreateMarketplaceDAO() (dao IWriterMarketplaceDAO, err error) {
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

	err = CheckMigrations(connection)
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	mpDao := &MarketplaceDAO{
		connection: connection,
	}

	return mpDao, nil
}

type MarketplaceDAO struct {
	connection *sql.DB
}

func (md MarketplaceDAO) Close() {
	md.connection.Close()
}
