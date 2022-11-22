package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type panicMessage struct {
	Message string `json:"message"`
}

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				var msg string
				for i := 2; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					msg += fmt.Sprintf("[%s:%d]\n", file, line)
				}
				log.Err(err).
					Str("uri", ctx.Request.RequestURI).
					Str("method", ctx.Request.Method).
					Str("stack", msg).
					Msgf("got http runtime panic")
				_ = ctx.Error(errors.WithStack(err))
				ctx.JSON(http.StatusInternalServerError, panicMessage{err.Error()})
			}
		}()
		ctx.Next()
	}
}
