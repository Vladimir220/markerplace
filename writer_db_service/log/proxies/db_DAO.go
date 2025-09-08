package proxies

import (
	"context"
	"fmt"
	"writer_db_service/db/DAO/postgres"
	"writer_db_service/env"
	"writer_db_service/models"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateDAOWithLog(ctx context.Context, dao postgres.IWriterMarketplaceDAO) postgres.IWriterMarketplaceDAO {
	return &DAOWithLog{
		original: dao,
		logger:   logger_lib.CreateLoggerGateway(ctx, "IMarketplaceDAO"),
		infoLogs: env.GetLogsConfig().PrintMarketplaceDAOInfo,
	}
}

type DAOWithLog struct {
	original postgres.IWriterMarketplaceDAO
	logger   logger_lib.ILogger
	infoLogs bool
}

func (mdwl *DAOWithLog) NewAnnouncement(announcement models.ExtendedAnnouncement) (err error) {
	logLabel := fmt.Sprintf("NewAnnouncement():[params:%v]:", announcement)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	err = mdwl.original.NewAnnouncement(announcement)
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

func (mdwl *DAOWithLog) UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (err error) {
	logLabel := fmt.Sprintf("UpdateAnnouncement():[params:%v]:", updatedAnnouncement)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	err = mdwl.original.UpdateAnnouncement(updatedAnnouncement)
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
