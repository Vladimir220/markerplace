package postgres

import (
	"context"
	"fmt"
	"marketplace/db/DAO/postgres/remote"
	"marketplace/models"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateDAOGateway(ctx context.Context) (dao IMarketplaceDAO, err error) {
	logLabel := "CreateDAOGateway():"
	logger := logger_lib.CreateLoggerGateway(ctx, "DAOGateway")

	grpcRemoteReader, err := remote.CreateGrpcReader(ctx)
	var remoteReaderUnavailable bool
	if err != nil {
		logger.WriteWarning(fmt.Sprintf("%s: %s: %v", logLabel, "remoteReader unavailable", err))
		remoteReaderUnavailable = true
	}

	kafkaRemoteWriter, err := remote.CreateKafkaWriter(ctx)
	var remoteWriterUnavailable bool
	if err != nil {
		logger.WriteWarning(fmt.Sprintf("%s: %s: %v", logLabel, "remoteWriter unavailable", err))
		remoteWriterUnavailable = true
	}

	localDao, err := CreateMarketplaceDAO()
	if err != nil {
		err = fmt.Errorf("%s: %v", logLabel, err)
		return
	}

	dao = &DAOGateway{
		remoteReaderUnavailable: remoteReaderUnavailable,
		remoteWriterUnavailable: remoteWriterUnavailable,
		localDAO:                localDao,
		remoteReader:            grpcRemoteReader,
		remoteWriter:            kafkaRemoteWriter,
		ctx:                     ctx,
		logger:                  logger,
	}

	return
}

// Proxy for IMarketplaceDAO
// Adapter for IReader and IWriter
type DAOGateway struct {
	remoteReaderUnavailable bool
	remoteWriterUnavailable bool
	localDAO                IMarketplaceDAO
	remoteReader            remote.IReader
	remoteWriter            remote.IWriter
	ctx                     context.Context
	logger                  logger_lib.ILogger
}

func (dao DAOGateway) GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error) {
	if !dao.remoteReaderUnavailable {
		announcement, err = dao.remoteReader.GetAnnouncements(orderType, minPrice, maxPrice, page)
		if err == nil {
			return
		}
	}

	dao.logger.WriteWarning(fmt.Sprintf("%s: %v", "GetAnnouncements(): remote remoteWriter unavailable", err))
	announcement, err = dao.localDAO.GetAnnouncements(orderType, minPrice, maxPrice, page)
	if err != nil {
		err = fmt.Errorf("%s: %v", "GetAnnouncements():", err)
		return
	}

	return
}

func (dao DAOGateway) NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error) {
	if !dao.remoteWriterUnavailable {
		err = dao.remoteWriter.NewAnnouncement(announcement)
		if err == nil {
			return
		}
	}

	dao.logger.WriteWarning(fmt.Sprintf("%s: %v", "NewAnnouncement(): remote remoteWriter unavailable", err))
	resAnnouncement, err = dao.localDAO.NewAnnouncement(announcement)
	if err != nil {
		err = fmt.Errorf("%s: %v", "NewAnnouncement():", err)
		return
	}

	return
}

func (dao DAOGateway) Registr(login, password string) (user models.User, isAlreadyExist bool, err error) {
	user, isAlreadyExist, err = dao.localDAO.Registr(login, password)
	if err != nil {
		err = fmt.Errorf("%s:%v", "Registr():", err)
	}
	return
}

func (dao DAOGateway) Close() {
	dao.localDAO.Close()
}

func (dao DAOGateway) GetUser(login string) (user models.User, password string, isFound bool, err error) {
	user, password, isFound, err = dao.localDAO.GetUser(login)
	if err != nil {
		err = fmt.Errorf("%s:%v", "GetUser():", err)
	}
	return
}

func (dao DAOGateway) UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error) {
	if !dao.remoteWriterUnavailable {
		err = dao.remoteWriter.UpdateAnnouncement(updatedAnnouncement)
		if err == nil {
			return
		}
	}

	dao.logger.WriteWarning(fmt.Sprintf("%s: %v", "UpdateAnnouncement(): remote remoteWriter unavailable", err))
	resAnnouncement, err = dao.localDAO.UpdateAnnouncement(updatedAnnouncement)
	if err != nil {
		err = fmt.Errorf("%s: %v", "UpdateAnnouncement():", err)
		return
	}

	return
}

func (dao DAOGateway) DeleteAnnouncement(announcementId uint) (err error) {
	if !dao.remoteWriterUnavailable {
		err = dao.remoteWriter.DeleteAnnouncement(announcementId)
		if err == nil {
			return
		}
	}

	dao.logger.WriteWarning(fmt.Sprintf("%s: %v", "DeleteAnnouncement(): remote remoteWriter unavailable", err))
	err = dao.localDAO.DeleteAnnouncement(announcementId)
	if err != nil {
		err = fmt.Errorf("%s: %v", "DeleteAnnouncement():", err)
		return
	}

	return
}

func (dao DAOGateway) GetAuthorLogin(announcementId uint) (login string, isAnnouncementFound bool, err error) {
	login, isAnnouncementFound, err = dao.localDAO.GetAuthorLogin(announcementId)
	if err != nil {
		err = fmt.Errorf("%s:%v", "GetAuthorLogin():", err)
	}
	return
}
