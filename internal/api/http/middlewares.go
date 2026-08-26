package http

import (
	"cidadon/internal/auth"
	apperrors "cidadon/internal/handler"
	"cidadon/internal/provider"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		status, resp := apperrors.FromError(c.Errors.Last().Err)
		c.JSON(status, resp)
	}
}

type AuthMiddleware struct {
	JwtProvider provider.JwtProvider
}

func NewAuthMiddleware(JwtProvider provider.JwtProvider) *AuthMiddleware {
	return &AuthMiddleware{
		JwtProvider: JwtProvider,
	}
}

func (a *AuthMiddleware) AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Request.Cookie("accessToken")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		subject, err := a.JwtProvider.GetSubject(accessToken.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		subjectAuth, err := auth.GetSubjectAuth(subject)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userId", subjectAuth.UserID)
		c.Set("userRole", subjectAuth.Role)
		c.Next()
	}
}
