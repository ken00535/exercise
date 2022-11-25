package http

import (
	"net/http"

	"shorten/internal/app/entity"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// httpErr represents an error that occurred while handling a request.
type httpErr struct {
	Message string `json:"message"`
}

// ErrorHandler responds error response according to given error.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := errors.Unwrap(c.Errors[0])
		var appErr entity.AppError
		for _, err := range c.Errors {
			if errors.As(err, &appErr) {
				var status int
				switch {
				case errors.Is(appErr, entity.ErrInvalidInput):
					status = http.StatusBadRequest
				case errors.Is(appErr, entity.ErrUnauthorized):
					status = http.StatusUnauthorized
				case errors.Is(appErr, entity.ErrPermissionDenied):
					status = http.StatusForbidden
				case errors.Is(appErr, entity.ErrResourceIsEmpty):
					status = http.StatusUnprocessableEntity
				case errors.Is(appErr, entity.ErrResourceHasExisted):
					status = http.StatusConflict
				case errors.Is(appErr, entity.ErrResourceNotFound):
					status = http.StatusNotFound
				case errors.Is(appErr, entity.ErrTransactionNotCompleted):
					status = http.StatusPreconditionRequired
				case errors.Is(appErr, entity.ErrRuntimePanic):
					status = http.StatusInternalServerError
				default:
					status = http.StatusInternalServerError
				}
				log.Err(errors.Unwrap(err)).Msgf("%s", err.Error())
				c.JSON(status, httpErr{Message: appErr.Error()})
				return
			}
		}
		log.Err(err).Msgf("%s", err.Error())
		c.JSON(http.StatusInternalServerError, httpErr{Message: err.Error()})
	}
}
