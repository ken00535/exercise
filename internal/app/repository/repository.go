//go:generate mockgen -source ../../entity/repository_db.go -destination mock/mock.go -package mock
package repository

import (
	"context"
	"errors"

	"shorten/internal/app/entity"
)

type repository struct {
	db    entity.RepositoryDb
	cache entity.RepositoryCache
	mem   entity.RepositoryMem
}

func New(db entity.RepositoryDb, cache entity.RepositoryCache, mem entity.RepositoryMem) entity.Repository {
	r := repository{
		db:    db,
		cache: cache,
		mem:   mem,
	}
	return &r
}

func (r *repository) SaveShortenUrl(ctx context.Context, url *entity.Url) (*entity.Url, error) {
	var err error
	url, err = r.db.SaveShortenUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	url, err = r.cache.SaveShortenUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	url, err = r.mem.SaveShortenUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (r *repository) GetUrl(ctx context.Context, path string) (*entity.Url, error) {
	var err error
	var url *entity.Url
	url, err = r.mem.GetUrl(ctx, path)
	if err == nil {
		return url, nil
	}
	if !errors.Is(err, entity.ErrResourceNotFound) {
		return nil, err
	}
	url, err = r.cache.GetUrl(ctx, path)
	if err == nil {
		return url, nil
	}
	if !errors.Is(err, entity.ErrResourceNotFound) {
		return nil, err
	}
	url, err = r.db.GetUrl(ctx, path)
	if err != nil {
		return nil, err
	}
	_, err = r.cache.SaveShortenUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	_, err = r.mem.SaveShortenUrl(ctx, url)
	if err != nil {
		return nil, err
	}
	return url, nil
}
