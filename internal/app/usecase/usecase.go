package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"shorten/internal/app/entity"

	"github.com/pkg/errors"
)

// Config setting log config
type Config struct {
	Scheme   string `yaml:"scheme" mapstructure:"scheme"`
	Hostname string `yaml:"hostname" mapstructure:"hostname"`
}

type usecase struct {
	db  entity.Repository
	cfg Config
}

func New(repo entity.Repository, cfg Config) entity.Usecase {
	u := usecase{
		db:  repo,
		cfg: cfg,
	}
	return &u
}

func (u *usecase) UploadUrl(ctx context.Context, url string, expireAt time.Time) (*entity.Url, error) {
	h := sha256.New()
	h.Write([]byte(url))
	bs := h.Sum(nil)
	enc := base64.URLEncoding.EncodeToString(bs)
	shorten := &entity.Url{
		Url:      url,
		ShortUrl: enc[:8],
		ExpireAt: expireAt,
	}
	shorten, err := u.db.SaveShortenUrl(ctx, shorten)
	if err != nil {
		return nil, err
	}
	shorten.ShortUrl = u.cfg.Scheme + "://" + u.cfg.Hostname + "/" + shorten.ShortUrl
	return shorten, nil
}

func (u *usecase) GetUrl(ctx context.Context, path string) (*entity.Url, error) {
	url, err := u.db.GetUrl(ctx, path)
	if err != nil {
		return nil, err
	}
	if url.ExpireAt.Before(time.Now()) {
		return nil, errors.Wrap(entity.ErrResourceNotFound, "over expire time")
	}
	return url, nil
}
