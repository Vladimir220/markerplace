package postgres

import "fmt"

func (md WriterMarketplaceDAO) DeleteAnnouncement(announcementId uint) (err error) {
	queryStr := `DELETE FROM announcements WHERE id=$1;`

	_, err = md.connection.Exec(queryStr, announcementId)
	if err != nil {
		err = fmt.Errorf("MarketplaceDAO:DeleteAnnouncement: %v", err)
		return
	}

	return
}
