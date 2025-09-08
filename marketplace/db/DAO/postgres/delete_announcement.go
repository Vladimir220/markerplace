package postgres

import (
	"database/sql"
	"fmt"
)

func (md MarketplaceDAO) DeleteAnnouncement(announcementId uint) (err error) {
	queryStr := `DELETE FROM announcements WHERE id=$1;`

	_, err = md.connection.Exec(queryStr, announcementId)
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("MarketplaceDAO:DeleteAnnouncement: %v", err)
		return
	}

	return
}
