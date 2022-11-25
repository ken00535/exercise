package entity

import "context"

type RepositoryCache interface {
	SaveShortenUrl(ctx context.Context, url *Url) (*Url, error)
	GetUrl(ctx context.Context, shortenUrl string) (*Url, error)
}
