//go:generate mockgen -source ../../entity/repository_mem.go -destination mock/mock.go -package mock
package mem

import (
	"context"
	"encoding/json"
	"shorten/internal/app/entity"

	"github.com/allegro/bigcache/v3"
	"github.com/pkg/errors"
)

type repository struct {
	cache *bigcache.BigCache
}

func New(mem *bigcache.BigCache) entity.RepositoryMem {
	r := repository{
		cache: mem,
	}
	return &r
}

func (r *repository) SaveShortenUrl(ctx context.Context, url *entity.Url) (*entity.Url, error) {
	data, err := json.Marshal(url)
	if err != nil {
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	if err := r.cache.Set(url.ShortUrl, data); err != nil {
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	return url, nil
}

func (r *repository) GetUrl(ctx context.Context, path string) (*entity.Url, error) {
	res, err := r.cache.Get(path)
	if err != nil {
		switch {
		case errors.Is(err, bigcache.ErrEntryNotFound):
			return nil, errors.Wrapf(entity.ErrResourceNotFound, err.Error())
		}
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	url := &entity.Url{}
	if err := json.Unmarshal([]byte(res), &url); err != nil {
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	return url, nil
}
