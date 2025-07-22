package models

import "time"

type User struct {
	Login, Group string
}

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

type Announcements struct {
	Ans         []ExtendedAnnouncement `json:"announcements"`
	Page        uint                   `json:"page"`
	MaxPage     uint                   `json:"maxPage"`
	MinPrice    uint                   `json:"minPrice"`
	PriceFilter bool                   `json:"priceFilter"`
}
