//go:generate mockgen -source ../../entity/repository_cache.go -destination mock/mock.go -package mock
package cache

import (
	"context"
	"encoding/json"
	"shorten/internal/app/entity"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type repository struct {
	cache *redis.Client
}

func New(db *redis.Client) entity.RepositoryCache {
	r := repository{
		cache: db,
	}
	return &r
}

func (r *repository) SaveShortenUrl(ctx context.Context, url *entity.Url) (*entity.Url, error) {
	dao := &daoUrl{}
	_ = copier.Copy(&dao, &url)
	data, err := json.Marshal(url)
	if err != nil {
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	if err := r.cache.Set(ctx, url.ShortUrl, data, time.Hour).Err(); err != nil {
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	url.ID = dao.ShortUrl
	return url, nil
}

func (r *repository) GetUrl(ctx context.Context, path string) (*entity.Url, error) {
	dao := &daoUrl{}
	res, err := r.cache.Get(ctx, path).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, errors.Wrapf(entity.ErrResourceNotFound, err.Error())
		}
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	if err := json.Unmarshal([]byte(res), &dao); err != nil {
		return nil, errors.Wrapf(entity.ErrInternal, err.Error())
	}
	url := &entity.Url{}
	_ = copier.Copy(&url, &dao)
	return url, nil
}
