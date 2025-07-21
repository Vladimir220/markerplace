package postgres

import (
	"main/models"
	"main/tools/db"
)

type IMarcketplaceDAO interface {
	GetAnnouncements(orderType *string, minPrice, maxPrice *uint, offset, limit uint) (announcement models.ExtendedAnnouncement, err error)
	GetUser(login string) (user models.User, password string, isFound bool, err error)
	NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error)
	Registr(login, password string) (user models.User, err error)
	Close()
}

func CreateMarcketplaceDAO() IMarcketplaceDAO {
	return MarcketplaceDAO{
		сonnectionPool: db.CreateConnectionPool(),
	}
}

type MarcketplaceDAO struct {
	сonnectionPool db.IConnectionPool
}

func (md MarcketplaceDAO) Close() {
	md.сonnectionPool.Close()
}
