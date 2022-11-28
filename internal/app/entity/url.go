package entity

import "time"

type Url struct {
	ID       string
	Url      string
	ShortUrl string
	ExpireAt time.Time
}
