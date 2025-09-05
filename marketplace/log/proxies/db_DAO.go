package proxies

import (
	"context"
	"fmt"
	"marketplace/db/DAO/postgres"
	"marketplace/log"
	"marketplace/models"
)

func CreateDAOWithLog(ctx context.Context, dao postgres.IMarketplaceDAO, infoLogs bool) postgres.IMarketplaceDAO {
	return &DAOWithLog{
		original: dao,
		logger:   log.CreateLoggerAdapter(ctx, "IMarketplaceDAO"),
		infoLogs: infoLogs,
	}
}

type DAOWithLog struct {
	original postgres.IMarketplaceDAO
	logger   log.ILogger
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
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "get"))
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
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "get"))
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
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "get"))
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
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "get"))
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
