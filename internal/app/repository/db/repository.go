//go:generate mockgen -source ../../entity/repository_db.go -destination mock/mock.go -package mock
package db

import (
	"assignment/internal/app/entity"
)

type repository struct{}

func New() entity.RepositoryDb {
	r := repository{}
	return &r
}
