package usecase

import (
	"assignment/internal/app/entity"
)

type usecase struct {
}

func New() entity.Usecase {
	u := usecase{}
	return &u
}
