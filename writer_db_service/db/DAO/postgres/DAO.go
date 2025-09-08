package postgres

import (
	"database/sql"
	"fmt"
	"time"
	"writer_db_service/models"

	_ "github.com/lib/pq"
)

type IWriterMarketplaceDAO interface {
	NewAnnouncement(announcement models.ExtendedAnnouncement) (err error)
	UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (err error)
	DeleteAnnouncement(announcementId uint) (err error)
	Close()
}

func CreateWriterMarketplaceDAO() (dao IWriterMarketplaceDAO, err error) {
	logLabel := "CreateMarketplaceDAO():"

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

	mpDao := &WriterMarketplaceDAO{
		connection: connection,
	}

	return mpDao, nil
}

type WriterMarketplaceDAO struct {
	connection *sql.DB
}

func (md WriterMarketplaceDAO) Close() {
	md.connection.Close()
}
