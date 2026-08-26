package handler

import (
	"cidadon/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *Handler {
	return &Handler{
		AuthService: authService,
	}
}

func (ah *Handler) Login(c *gin.Context) {
	var loginRequest service.LoginInput
	if err := c.ShouldBind(&loginRequest); err != nil {
		c.Error(err)
		return
	}

	session, err := ah.AuthService.Login(c.Request.Context(), loginRequest)
	if err != nil {
		c.Error(err)
		return
	}

	ah.setAuthCookies(c, session)
	c.Status(http.StatusCreated)
}

func (ah *Handler) RegisterCitizen(c *gin.Context) {
	var registerInput service.RegisterCitizenInput
	if err := c.ShouldBind(&registerInput); err != nil {
		c.Error(err)
		return
	}

	err := ah.AuthService.RegisterCitizen(c.Request.Context(), registerInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (ah *Handler) RegisterCouncillor(c *gin.Context) {
	var registerInput service.RegisterCouncillorInput
	if err := c.ShouldBind(&registerInput); err != nil {
		c.Error(err)
		return
	}

	err := ah.AuthService.RegisterCouncillor(c.Request.Context(), registerInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (ah *Handler) RegisterOfficeMember(c *gin.Context) {
	var registerInput service.RegisterOfficeMemberInput
	if err := c.ShouldBind(&registerInput); err != nil {
		c.Error(err)
		return
	}

	err := ah.AuthService.RegisterOfficeMember(c.Request.Context(), registerInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusCreated)
}

func (ah *Handler) setAuthCookies(c *gin.Context, session service.LoginOutput) {
	now := time.Now()
	c.SetCookie(accessTokenCookieKey, session.AccessToken, int(session.AccessTokenExpiresIn.Sub(now).Seconds()), "/", "", false, true)
	c.SetCookie(refreshTokenCookieKey, session.RefreshToken, int(session.RefreshTokenExpiresIn.Sub(now).Seconds()), "/", "", false, true)
}
