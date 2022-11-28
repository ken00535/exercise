package cache

import (
	"time"
)

type daoUrl struct {
	ID       string
	Url      string
	ShortUrl string
	ExpireAt time.Time
}
