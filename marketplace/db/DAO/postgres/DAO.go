package postgres

import (
	"database/sql"
	"fmt"
	"marketplace/models"
	"time"

	_ "github.com/lib/pq"
)

type IMarketplaceDAO interface {
	GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error)
	GetUser(login string) (user models.User, password string, isFound bool, err error)
	NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error)
	Registr(login, password string) (user models.User, isAlreadyExist bool, err error)
	UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error)
	DeleteAnnouncement(announcementId uint) (err error)
	GetAuthorLogin(announcementId uint) (login string, isAnnouncementFound bool, err error)
	Close()
}

func CreateMarketplaceDAO() (dao IMarketplaceDAO, err error) {
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
