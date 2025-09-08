package postgres

import (
	"fmt"
	"marketplace/models"
)

func (md MarketplaceDAO) UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error) {
	queryStr := `UPDATE announcements SET
				title=$1,
				body=$2,
				pic_link=$3,
				price=$4,
				date=$5,
				author_login=$6
				WHERE id = $7
				RETURNING title, body, pic_link, price, date, author_login, id;
	`

	err = md.connection.QueryRow(queryStr, updatedAnnouncement.An.Title, updatedAnnouncement.An.Body, updatedAnnouncement.An.PicLink, updatedAnnouncement.An.Price, updatedAnnouncement.Date, updatedAnnouncement.AuthorLogin, updatedAnnouncement.Id).
		Scan(&resAnnouncement.An.Title, &resAnnouncement.An.Body, &resAnnouncement.An.PicLink, &resAnnouncement.An.Price, &resAnnouncement.Date, &resAnnouncement.AuthorLogin, &resAnnouncement.Id)
	if err != nil {
		err = fmt.Errorf("MarketplaceDAO:UpdateAnnouncement: %v", err)
		return
	}

	return
}
