package postgres

import (
	"errors"
	"fmt"
	"writer_db_service/models"
)

func (md WriterMarketplaceDAO) NewAnnouncement(announcement models.ExtendedAnnouncement) (err error) {
	if announcement.An.Title == "" || announcement.An.Body == "" || announcement.An.PicLink == "" || announcement.AuthorLogin == "" {
		err = errors.New("WriterMarketplaceDAO:NewAnnouncement: Title or Body or PicLink or AuthorLogin not specified")
		return
	}

	queryStr := `INSERT INTO announcements (title, body, pic_link, price, date, author_login) 
				VALUES ($1, $2, $3, $4, $5, $6);`

	_, err = md.connection.Exec(queryStr, announcement.An.Title, announcement.An.Body, announcement.An.PicLink, announcement.An.Price, announcement.Date, announcement.AuthorLogin)
	if err != nil {
		err = fmt.Errorf("WriterMarketplaceDAO:NewAnnouncement(): %v", err)
		return
	}

	return
}
