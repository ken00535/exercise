package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"shorten/internal/app/entity"
)

// Config setting log config
type Config struct {
	Scheme   string `yaml:"scheme" mapstructure:"scheme"`
	Hostname string `yaml:"hostname" mapstructure:"hostname"`
}

type usecase struct {
	db  entity.RepositoryDb
	cfg Config
}

func New(repo entity.RepositoryDb, cfg Config) entity.Usecase {
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

func (u *usecase) GetUrl(ctx context.Context, url string) (*entity.Url, error) {
	shorten := &entity.Url{
		Url: "https://google.com",
	}
	return shorten, nil
}
