package http

import (
	"fmt"
	"net/http"
	"time"

	"shorten/internal/app/entity"
	"shorten/internal/infra/clock"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Delivery struct {
	usecase entity.Usecase
	clock   clock.Clock
}

func New(e *gin.Engine, u entity.Usecase) *Delivery {
	d := &Delivery{
		usecase: u,
		clock:   clock.Real{},
	}
	e.Use(ErrorHandler())
	e.GET("/:url", d.ServeShortenUrl)
	apiV1 := e.Group("/api/v1")
	apiV1.POST("/urls", d.UploadUrl)
	return d
}

func (d *Delivery) ServeShortenUrl(c *gin.Context) {
	source := c.Param("url")
	if source == "" {
		_ = c.Error(errors.Wrap(entity.ErrInvalidInput, "no url"))
		return
	}
	url, err := d.usecase.GetUrl(c, source)
	if err != nil {
		_ = c.Error(err)
		return
	}
	expirePeriod := url.ExpireAt.Sub(d.clock.Now())
	c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(expirePeriod.Seconds())))
	c.Redirect(http.StatusFound, url.Url)
}

func (d *Delivery) UploadUrl(c *gin.Context) {
	type Body struct {
		Url      string    `json:"url" binding:"required,url"`
		ExpireAt time.Time `json:"expireAt" binding:"required"`
	}
	req := Body{}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errors.Wrap(entity.ErrInvalidInput, err.Error()))
		return
	}
	url, err := d.usecase.UploadUrl(c, req.Url, req.ExpireAt)
	if err != nil {
		_ = c.Error(err)
		return
	}
	type Response struct {
		ID       int64  `json:"id"`
		ShortUrl string `json:"shortUrl"`
	}
	c.JSON(http.StatusOK, Response{
		ID:       url.ID,
		ShortUrl: url.ShortUrl,
	})
}
