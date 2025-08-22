package postgres

import (
	"main/db/tools"
	"main/models"
)

type IMarketplaceDAO interface {
	GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error)
	GetUser(login string) (user models.User, password string, isFound bool, err error)
	NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error)
	Registr(login, password string) (user models.User, isAlreadyExist bool, err error)
	Close()
}

func CreateMarketplaceDAO() IMarketplaceDAO {
	return MarketplaceDAO{
		сonnectionPool: tools.CreateConnectionPool(),
	}
}

type MarketplaceDAO struct {
	сonnectionPool tools.IConnectionPool
}

func (md MarketplaceDAO) Close() {
	md.сonnectionPool.Close()
}
