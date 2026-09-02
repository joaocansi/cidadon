package http

import (
	"cidadon/internal/adapters/external/provider"
	service "cidadon/internal/application/contract"
	"cidadon/internal/domain/entity"

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
			c.Error(service.Unauthorized("authentication required"))
			c.Abort()
			return
		}

		subject, err := a.JwtProvider.GetSubject(accessToken.Value)
		if err != nil {
			c.Error(service.Unauthorized("invalid session"))
			c.Abort()
			return
		}

		subjectAuth := entity.AuthSubject{}
		err = subjectAuth.FromString(subject)

		if err != nil {
			c.Error(service.Unauthorized("invalid session"))
			c.Abort()
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

		c.Error(service.Forbidden("role is not allowed"))
		c.Abort()
		return
	}
}

// OptionalAuthHandler enriches the request with the authenticated subject when
// a valid session exists. It deliberately treats a missing or expired cookie as
// an anonymous request, which is useful for session discovery endpoints.
func (a *AuthMiddleware) OptionalAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Request.Cookie("accessToken")
		if err != nil {
			c.Next()
			return
		}

		subject, err := a.JwtProvider.GetSubject(accessToken.Value)
		if err != nil {
			c.Next()
			return
		}

		subjectAuth := entity.AuthSubject{}
		if err := subjectAuth.FromString(subject); err == nil {
			c.Set("userId", subjectAuth.UserID)
			c.Set("userRole", subjectAuth.Role)
		}
		c.Next()
	}
}
