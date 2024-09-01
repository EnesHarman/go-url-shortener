package model

import "time"

type ClickEvent struct {
	UserId string    `json:"user_id"`
	Ts     time.Time `json:"ts"`
	UrlId  int       `json:"url_id"`
}
