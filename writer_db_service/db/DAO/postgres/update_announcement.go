package postgres

import (
	"fmt"
	"writer_db_service/models"
)

func (md WriterMarketplaceDAO) UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (err error) {
	queryStr := `UPDATE announcements SET
				title=$1,
				body=$2,
				pic_link=$3,
				price=$4,
				date=$5,
				author_login=$6
				WHERE id = $7;
	`

	_, err = md.connection.Exec(queryStr, updatedAnnouncement.An.Title, updatedAnnouncement.An.Body, updatedAnnouncement.An.PicLink, updatedAnnouncement.An.Price, updatedAnnouncement.Date, updatedAnnouncement.AuthorLogin, updatedAnnouncement.Id)
	if err != nil {
		err = fmt.Errorf("WriterMarketplaceDAO:UpdateAnnouncement: %v", err)
		return
	}

	return
}
