package postgres

import (
	"errors"
	"fmt"
	"main/models"
)

func (md MarcketplaceDAO) NewAnnouncement(announcement models.ExtendedAnnouncement) (resAnnouncement models.ExtendedAnnouncement, err error) {
	if announcement.An.Title == "" || announcement.An.Body == "" || announcement.An.PicLink == "" || announcement.AuthorLogin == "" {
		err = errors.New("MarcketplaceDAO:NewAnnouncement: Title or Body or PicLink or AuthorLogin not specified")
		return
	}

	queryStr := `INSERT INTO announcements (title, body, pic_link, price, date, author_login) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING title, body, pic_link, price, date, author_login, id;`

	connection := md.сonnectionPool.GetConnection()

	err = connection.QueryRow(queryStr, announcement.An.Title, announcement.An.Body, announcement.An.PicLink, announcement.An.Price, announcement.Date, announcement.AuthorLogin).
		Scan(&resAnnouncement.An.Title, &resAnnouncement.An.Body, &resAnnouncement.An.PicLink, &resAnnouncement.An.Price, &resAnnouncement.Date, &resAnnouncement.AuthorLogin, resAnnouncement.Id)
	if err != nil {
		err = fmt.Errorf("MarcketplaceDAO:NewAnnouncement: %v", err)
		return
	}

	return
}
