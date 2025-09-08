package proxies

import (
	"context"
	"fmt"
	"reader_db_service/db/DAO/postgres"
	"reader_db_service/env"
	"reader_db_service/models"

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

func (mdwl *DAOWithLog) Close() {
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo("Close()")
	}
	mdwl.original.Close()
}
