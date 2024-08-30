package model

import "time"

type Url struct {
	Id         int       `json:"_id"`
	RealUrl    string    `json:"real_url"`
	ShortUrl   string    `json:"short_url"`
	Ts         time.Time `json:"ts"`
	ExpireDate time.Time `json:"expire_date"`
}
