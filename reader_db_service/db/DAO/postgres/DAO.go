package postgres

import (
	"database/sql"
	"fmt"
	"reader_db_service/models"
	"time"

	_ "github.com/lib/pq"
)

type IReaderMarketplaceDAO interface {
	GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error)
	Close()
}

func CreateReaderMarketplaceDAO() (dao IReaderMarketplaceDAO, err error) {
	logLabel := "CreateReaderMarketplaceDAO():"

	err = checkDbExistence(5, time.Second)
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	connection, err := connect()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	mpDao := &ReaderMarketplaceDAO{
		connection: connection,
	}

	return mpDao, nil
}

type ReaderMarketplaceDAO struct {
	connection *sql.DB
}

func (md ReaderMarketplaceDAO) Close() {
	md.connection.Close()
}
