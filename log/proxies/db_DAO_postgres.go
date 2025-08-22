package proxies

import (
	"fmt"
	"main/db/DAO/postgres"
	"main/log"
	"main/models"
)

func CreateDAOWithLog(dao postgres.IMarketplaceDAO) postgres.IMarketplaceDAO {
	return DAOWithLog{
		original: dao,
		logger:   log.CreateLogger("IMarketplaceDAO"),
	}
}

type DAOWithLog struct {
	original postgres.IMarketplaceDAO
	logger   log.ILogger
}

func (mdwl DAOWithLog) GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error) {
	logLabel := "GetAnnouncements():"
	announcement, err = mdwl.original.GetAnnouncements(orderType, minPrice, maxPrice, page)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl DAOWithLog) GetUser(login string) (user models.User, password string, isFound bool, err error) {
	logLabel := "GetUser():"
	user, password, isFound, err = mdwl.original.GetUser(login)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl DAOWithLog) NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error) {
	logLabel := "NewAnnouncement():"
	resAnnouncement, err = mdwl.original.NewAnnouncement(announcement)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl DAOWithLog) Registr(login, password string) (user models.User, isAlreadyExist bool, err error) {
	logLabel := "Registr():"
	user, isAlreadyExist, err = mdwl.original.Registr(login, password)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl DAOWithLog) Close() {
	mdwl.original.Close()
}
