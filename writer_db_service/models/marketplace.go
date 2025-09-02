package models

import "time"

type Announcement struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	PicLink string `json:"pic_link"`
	Price   uint   `json:"price"`
}

type ExtendedAnnouncement struct {
	An          Announcement `json:"announcement"`
	Id          uint         `json:"id"`
	AuthorLogin string       `json:"author-login"`
	Date        time.Time    `json:"date"`
	Yours       bool         `json:"yours"`
}
