package main

import (
	"context"
	"fmt"
	"reader_db_service/db/DAO/postgres"
	"reader_db_service/gen"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateServer(ctx context.Context, dao postgres.IReaderMarketplaceDAO) gen.ReaderServer {
	return &Server{
		dao:    dao,
		logger: logger_lib.CreateLoggerGateway(ctx, "Server"),
	}
}

type Server struct {
	gen.UnimplementedReaderServer
	dao    postgres.IReaderMarketplaceDAO
	logger logger_lib.ILogger
}

func (s Server) GetAnnouncements(ctx context.Context, req *gen.AnnouncementsRequest) (resp *gen.AnnouncementsResponse, err error) {
	logLabel := fmt.Sprintf("%s:%s:", serviceName, "GetAnnouncements()")

	if req == nil {
		err = fmt.Errorf("%s %s", logLabel, "AnnouncementsRequest is nil")
		s.logger.WriteError(err.Error())
		return
	}

	var minPrice, maxPrice *uint
	if req.MinPrice != nil {
		minPrice = new(uint)
		*minPrice = uint(*req.MinPrice)
	}
	if req.MaxPrice != nil {
		maxPrice = new(uint)
		*maxPrice = uint(*req.MaxPrice)
	}

	announcement, err := s.dao.GetAnnouncements(req.OrderType, minPrice, maxPrice, uint(req.Page))

	ans := make([]*gen.Announcement, len(announcement.Ans))
	for i := range announcement.Ans {
		ans[i] = &gen.Announcement{}
		ans[i].Title = announcement.Ans[i].An.Title
		ans[i].Body = announcement.Ans[i].An.Body
		ans[i].PicLink = announcement.Ans[i].An.PicLink
		ans[i].Price = uint64(announcement.Ans[i].An.Price)
		ans[i].Id = uint64(announcement.Ans[i].Id)
		ans[i].AuthorLogin = announcement.Ans[i].AuthorLogin
		ans[i].DateUnixTimestamp = announcement.Ans[i].Date.Unix()
	}

	resp = &gen.AnnouncementsResponse{
		Ans:         ans,
		Page:        uint64(announcement.Page),
		MaxPage:     uint64(announcement.MaxPage),
		MinPrice:    uint64(announcement.MinPrice),
		PriceFilter: announcement.PriceFilter,
	}
	return
}
