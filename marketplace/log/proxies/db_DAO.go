package proxies

import (
	"context"
	"fmt"
	"marketplace/db/DAO/postgres"
	"marketplace/env"
	"marketplace/models"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateDAOWithLog(ctx context.Context, dao postgres.IMarketplaceDAO) postgres.IMarketplaceDAO {
	return &DAOWithLog{
		original: dao,
		logger:   logger_lib.CreateLoggerGateway(ctx, "IMarketplaceDAO"),
		infoLogs: env.GetLogsConfig().PrintMarketplaceDAOInfo,
	}
}

type DAOWithLog struct {
	original postgres.IMarketplaceDAO
	logger   logger_lib.ILogger
	infoLogs bool
}

func (mdwl *DAOWithLog) GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error) {
	var orderTypeStr, minPriceStr, maxPriceStr string
	if orderType == nil {
		orderTypeStr = "nil"
	} else {
		orderTypeStr = *orderType
	}
	if minPrice == nil {
		minPriceStr = "nil"
	} else {
		minPriceStr = fmt.Sprint(*minPrice)
	}
	if maxPrice == nil {
		maxPriceStr = "nil"
	} else {
		maxPriceStr = fmt.Sprint(*maxPrice)
	}
	logLabel := fmt.Sprintf("GetAnnouncements():[params:%s,%s,%s,%d]:", orderTypeStr, minPriceStr, maxPriceStr, page)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}

	announcement, err = mdwl.original.GetAnnouncements(orderType, minPrice, maxPrice, page)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) GetUser(login string) (user models.User, password string, isFound bool, err error) {
	logLabel := fmt.Sprintf("GetUser():[params:%s]:", login)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	user, password, isFound, err = mdwl.original.GetUser(login)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error) {
	logLabel := fmt.Sprintf("NewAnnouncement():[params:%v]:", announcement)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	resAnnouncement, err = mdwl.original.NewAnnouncement(announcement)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) Registr(login, password string) (user models.User, isAlreadyExist bool, err error) {
	logLabel := fmt.Sprintf("Registr():[params:%s,%s]:", login, "***")
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	user, isAlreadyExist, err = mdwl.original.Registr(login, password)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) Close() {
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo("Close()")
	}
	mdwl.original.Close()
}

func (mdwl *DAOWithLog) UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error) {
	logLabel := fmt.Sprintf("UpdateAnnouncement():[params:%v]:", updatedAnnouncement)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	resAnnouncement, err = mdwl.original.UpdateAnnouncement(updatedAnnouncement)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) DeleteAnnouncement(announcementId uint) (err error) {
	logLabel := fmt.Sprintf("DeleteAnnouncement():[params:%v]:", announcementId)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	err = mdwl.original.DeleteAnnouncement(announcementId)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) GetAuthorLogin(announcementId uint) (login string, isAnnouncementFound bool, err error) {
	logLabel := fmt.Sprintf("GetAuthorLogin():[params:%d]:", announcementId)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	login, isAnnouncementFound, err = mdwl.original.GetAuthorLogin(announcementId)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}
