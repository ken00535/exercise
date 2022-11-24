package entity

import (
	"context"
	"time"
)

type Usecase interface {
	UploadUrl(ctx context.Context, url string, expireAt time.Time) (*Url, error)
	GetUrl(ctx context.Context, url string) (*Url, error)
}
