package remote

import (
	"context"
	"fmt"
	"marketplace/env"
	"marketplace/gen/reader_db_service"
	"marketplace/models"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IReader interface {
	GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error)
}

func CreateGrpcReader(ctx context.Context) (reader IReader, err error) {
	data, err := env.GetServicesData()
	if err != nil {
		err = fmt.Errorf("CreateGrpcReader(): %v", err)
		return
	}

	reader = &GrpcReader{
		host: data.ReaderGRPCHost,
		ctx:  ctx,
	}

	return
}

type GrpcReader struct {
	host string
	ctx  context.Context
}

func (gr GrpcReader) GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcement models.Announcements, err error) {
	logLabel := "GrpcReader:GetAnnouncements():"

	conn, err := grpc.NewClient(gr.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	defer conn.Close()

	client := reader_db_service.NewReaderClient(conn)

	var minPrice64, maxPrice64 *uint64
	if minPrice != nil {
		minPrice64 = new(uint64)
		*minPrice64 = uint64(*minPrice)
	}
	if maxPrice != nil {
		maxPrice64 = new(uint64)
		*maxPrice64 = uint64(*maxPrice)
	}

	params := &reader_db_service.AnnouncementsRequest{
		OrderType: orderType,
		MinPrice:  minPrice64,
		MaxPrice:  maxPrice64,
		Page:      uint64(page),
	}

	resp, err := client.GetAnnouncements(gr.ctx, params)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	announcement.Ans = make([]models.ExtendedAnnouncement, len(resp.GetAns()))
	for i, v := range resp.GetAns() {
		announcement.Ans[i].An.Body = v.Body
		announcement.Ans[i].An.PicLink = v.PicLink
		announcement.Ans[i].An.Price = uint(v.Price)
		announcement.Ans[i].An.Title = v.Title
		announcement.Ans[i].AuthorLogin = v.AuthorLogin
		announcement.Ans[i].Date = time.Unix(v.DateUnixTimestamp, 0)
		announcement.Ans[i].Id = uint(v.Id)
	}
	announcement.MaxPage = uint(resp.GetMaxPage())
	announcement.MinPrice = uint(resp.GetMinPrice())
	announcement.Page = uint(resp.GetPage())
	announcement.PriceFilter = resp.GetPriceFilter()

	return
}
