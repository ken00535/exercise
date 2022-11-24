//go:generate mockgen -source ../../entity/repository_db.go -destination mock/mock.go -package mock
package db

import (
	"context"
	"shorten/internal/app/entity"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) entity.RepositoryDb {
	r := repository{
		db: db,
	}
	return &r
}

func (r *repository) SaveShortenUrl(ctx context.Context, url *entity.Url) (*entity.Url, error) {
	dao := &daoUrl{}
	_ = copier.Copy(&dao, &url)
	if res := r.db.Save(&dao); res.Error != nil {
		return nil, errors.Wrapf(entity.ErrInternal, res.Error.Error())
	}
	url.ID = dao.ID
	return url, nil
}
