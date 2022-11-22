package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"assignment/internal/app/entity"
)

type usecase struct {
}

func New() entity.Usecase {
	u := usecase{}
	return &u
}

func (u *usecase) UploadUrl(ctx context.Context, url string, expireAt time.Time) (*entity.Url, error) {
	h := sha256.New()
	h.Write([]byte(url))
	bs := h.Sum(nil)
	enc := base64.URLEncoding.EncodeToString(bs)
	shorten := &entity.Url{
		ID:       0,
		ShortUrl: enc[:8],
	}
	return shorten, nil
}
