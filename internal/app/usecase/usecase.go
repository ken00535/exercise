package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"shorten/internal/app/entity"
)

type usecase struct {
	repo entity.RepositoryDb
}

func New(repo entity.RepositoryDb) entity.Usecase {
	u := usecase{
		repo: repo,
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
	shorten, err := u.repo.SaveShortenUrl(ctx, shorten)
	if err != nil {
		return nil, err
	}
	return shorten, nil
}
