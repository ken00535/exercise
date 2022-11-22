package http

import (
	"assignment/internal/app/entity"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	usecase entity.Usecase
}

func New(e *gin.Engine, u entity.Usecase) *Delivery {
	d := &Delivery{
		usecase: u,
	}
	return d
}
