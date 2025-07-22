package postgres

import (
	"fmt"
	"main/models"
)

const pageSize = 10

func (md MarcketplaceDAO) GetAnnouncements(orderType *string, minPrice, maxPrice *uint, page uint) (announcements models.Announcements, err error) {
	announcements.Ans = make([]models.ExtendedAnnouncement, pageSize)
	announcements.PriceFilter = true
	var filter string
	if maxPrice != nil && minPrice != nil {
		filter = fmt.Sprintf("WHERE price >= %d AND price < %d", *minPrice, *maxPrice)
		announcements.MaxPage = *maxPrice
		announcements.MinPrice = *minPrice
	} else if maxPrice != nil {
		filter = fmt.Sprintf("WHERE price < %d", *maxPrice)
		announcements.MaxPage = *maxPrice
	} else if minPrice != nil {
		filter = fmt.Sprintf("WHERE price >= %d", *minPrice)
		announcements.MinPrice = *minPrice
	} else {
		announcements.PriceFilter = false
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

	if page == 0 {
		page = 1
	}
	offset := fmt.Sprintf("OFFSET %d", (page-1)*pageSize)
	limit := fmt.Sprintf("LIMIT %d", page*pageSize)

	queryStr := fmt.Sprintf("SELECT title, body, pic_link, price, date, author_login, id FROM announcements %s %s %s %s;", filter, order, offset, limit)

	connection := md.сonnectionPool.GetConnection()

	rows, err := connection.Query(queryStr)
	if err != nil {
		err = fmt.Errorf("MarcketplaceDAO:GetUser: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		announcement := models.ExtendedAnnouncement{}
		err = rows.Scan(&announcement.An.Title, &announcement.An.Body, &announcement.An.PicLink, &announcement.An.Price, &announcement.Date, &announcement.AuthorLogin, &announcement.Id)
		if err != nil {
			err = fmt.Errorf("MarcketplaceDAO:GetUser: %v", err)
			return
		}
		announcements.Ans = append(announcements.Ans, announcement)
	}

	return
}
