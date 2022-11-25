//go:generate mockgen -source ../../entity/repository_db.go -destination mock/mock.go -package mock
package db

import (
	"context"
	"shorten/internal/app/entity"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	if res := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "short_url"}},
				DoUpdates: clause.AssignmentColumns([]string{"url"}),
			},
		).
		Save(&dao); res.Error != nil {
		return nil, errors.Wrapf(entity.ErrInternal, res.Error.Error())
	}
	url.ID = dao.ID
	return url, nil
}

func (r *repository) GetUrl(ctx context.Context, path string) (*entity.Url, error) {
	dao := &daoUrl{}
	if res := r.db.WithContext(ctx).Where("short_url = ?", path).Take(&dao); res.Error != nil {
		switch {
		case errors.Is(res.Error, gorm.ErrRecordNotFound):
			return nil, errors.Wrapf(entity.ErrResourceNotFound, res.Error.Error())
		default:
			return nil, errors.Wrapf(entity.ErrInternal, res.Error.Error())
		}
	}
	url := &entity.Url{}
	_ = copier.Copy(&url, &dao)
	return url, nil
}
