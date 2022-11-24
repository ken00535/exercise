package db

import (
	"time"
)

type daoUrl struct {
	ID       int64
	Url      string
	ShortUrl string
	ExpireAt time.Time
}

func (d *daoUrl) TableName() string {
	return "url"
}
