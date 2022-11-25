package entity

import "context"

type Repository interface {
	SaveShortenUrl(ctx context.Context, url *Url) (*Url, error)
	GetUrl(ctx context.Context, shortenUrl string) (*Url, error)
}
