package entity

import "time"

type Url struct {
	ID       int64
	Url      string
	ShortUrl string
	ExpireAt time.Time
}
