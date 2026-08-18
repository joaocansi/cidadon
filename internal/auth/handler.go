package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuthService Service
}

func NewAuthHandler(authService Service) *Handler {
	return &Handler{
		AuthService: authService,
	}
}

func (ah *Handler) Login(c *gin.Context) {
	var loginRequest *LoginInput
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.Error(err)
		return
	}

	session, err := ah.AuthService.Login(c.Request.Context(), LoginInput{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	})

	if err != nil {
		c.Error(err)
		return
	}

	now := time.Now()
	c.Status(http.StatusCreated)
	c.SetCookie("accessToken", session.AccessToken, int(session.AccessTokenExpiresIn.Sub(now).Seconds()), "/", "", false, true)
	c.SetCookie("refreshToken", session.RefreshToken, int(session.RefreshTokenExpiresIn.Sub(now).Seconds()), "/", "", false, true)
}

func (ah *Handler) Register(c *gin.Context) {
	var registerInput RegisterInput
	if err := c.ShouldBind(&registerInput); err != nil {
		c.Error(err)
		return
	}

	err := ah.AuthService.Register(c.Request.Context(), registerInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

//func (ah *Handler) Profile(c *gin.Context) {
//	userID, ok := c.Get("userID")
//	if !ok {
//
//	}
//}
