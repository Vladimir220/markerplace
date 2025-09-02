package main

import (
	"context"
	"errors"
	"reader_db_service/DAO/postgres"
	"reader_db_service/gen"
	"strconv"
)

func CreateServer(dao postgres.IMarketplaceDAO) gen.ReaderServer {
	return &Server{
		dao: dao,
	}
}

type Server struct {
	gen.UnimplementedReaderServer
	dao postgres.IMarketplaceDAO
}

func (s Server) GetAnnouncements(ctx context.Context, req *gen.AnnouncementsRequest) (resp *gen.AnnouncementsResponse, err error) {
	if req == nil {
		err = errors.New("AnnouncementsRequest is nil")
		return
	}

	var minPrice, maxPrice *uint
	if req.MinPrice != nil {
		minPrice = new(uint)
		var buf uint64
		buf, err = strconv.ParseUint(*req.MinPrice, 10, 64)
		if err != nil {
			return
		}
		*minPrice = uint(buf)
	}
	if req.MaxPrice != nil {
		maxPrice = new(uint)
		var buf uint64
		buf, err = strconv.ParseUint(*req.MaxPrice, 10, 64)
		if err != nil {
			return
		}
		*maxPrice = uint(buf)
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
