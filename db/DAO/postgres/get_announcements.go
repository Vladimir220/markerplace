package postgres

import (
	"fmt"
	"main/models"
)

func (md MarcketplaceDAO) GetAnnouncements(orderType *string, minPrice, maxPrice *uint, offset, limit uint) (announcement models.ExtendedAnnouncement, err error) {
	var filter string
	if maxPrice != nil && minPrice != nil {
		filter = fmt.Sprintf("WHERE price >= %d AND price < %d", *minPrice, *maxPrice)
	} else if maxPrice != nil {
		filter = fmt.Sprintf("WHERE price < %d", *maxPrice)
	} else if minPrice != nil {
		filter = fmt.Sprintf("WHERE price >= %d", *minPrice)
	}

	order := "ORDER BY "
	if orderType != nil {
		switch *orderType {
		case "dasc":
			order = fmt.Sprintf("%s date ASC", order)
		case "ddesc":
			order = fmt.Sprintf("%s date DESC", order)
		case "pasc":
			order = fmt.Sprintf("%s price ASC", order)
		case "pdesc":
			order = fmt.Sprintf("%s price DESC", order)
		default:
			order = fmt.Sprintf("%s date ASC", order)
		}
	} else {
		order = fmt.Sprintf("%s date ASC", order)
	}

	queryStr := fmt.Sprintf("SELECT title, body, pic-link, price, date, author FROM announcements %s %s;", filter, order)

	connection := md.сonnectionPool.GetConnection()

	err = connection.QueryRow(queryStr).Scan(&announcement.An.Title, &announcement.An.Body, &announcement.An.PicLink, &announcement.An.Price, &announcement.Date, &announcement.AuthorLogin)
	if err != nil {
		err = fmt.Errorf("MarcketplaceDAO:GetUser: %v", err)
		return
	}

	return
}
