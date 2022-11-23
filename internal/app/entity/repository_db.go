package entity

import "context"

type RepositoryDb interface {
	SaveShortenUrl(ctx context.Context, url *Url) (*Url, error)
}
