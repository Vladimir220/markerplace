package postgres

import (
	"database/sql"
	"fmt"
)

func (md MarketplaceDAO) GetAuthorLogin(announcementId uint) (login string, isAnnouncementFound bool, err error) {
	queryStr := "SELECT author_login FROM announcements WHERE id=$1;"

	err = md.connection.QueryRow(queryStr, announcementId).Scan(&login)
	if err == sql.ErrNoRows {
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("MarketplaceDAO:GetAuthorLogin: %v", err)
		return
	} else {
		isAnnouncementFound = true
	}

	return
}
