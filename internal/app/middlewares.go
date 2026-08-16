package app

import (
	"cidadon/internal/app/shared/apperror"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		appErr := apperror.FromError(c.Errors.Last().Err)
		c.JSON(appErr.StatusCode(), appErr)
	}
}
