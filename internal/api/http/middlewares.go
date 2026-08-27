package http

import (
	"cidadon/internal/domain/entity"
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

		status, resp := FromError(c.Errors.Last().Err)
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

func (a *AuthMiddleware) AuthHandler(allowedRoles ...entity.UserRole) gin.HandlerFunc {
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

		subjectAuth := entity.AuthSubject{}
		err = subjectAuth.FromString(subject)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		for _, role := range allowedRoles {
			if subjectAuth.Role == role {
				c.Set("userId", subjectAuth.UserID)
				c.Set("userRole", subjectAuth.Role)
				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
