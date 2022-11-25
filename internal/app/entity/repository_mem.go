package entity

import "context"

type RepositoryMem interface {
	SaveShortenUrl(ctx context.Context, url *Url) (*Url, error)
	GetUrl(ctx context.Context, shortenUrl string) (*Url, error)
}
