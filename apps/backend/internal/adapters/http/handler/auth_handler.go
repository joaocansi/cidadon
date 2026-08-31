package handler

import (
	service "cidadon/internal/application/contract"
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
	if !bindRequest(c, &loginRequest) {
		return
	}

	session, err := ah.AuthService.Login(c.Request.Context(), loginRequest)
	if err != nil {
		c.Error(err)
		return
	}

	ah.setAuthCookies(c, session)
	c.JSON(http.StatusCreated, gin.H{
		"access_token_expires_in":  session.AccessTokenExpiresIn,
		"refresh_token_expires_in": session.RefreshTokenExpiresIn,
		"role":                     session.Role,
	})
}

func (ah *Handler) RegisterCitizen(c *gin.Context) {
	var registerInput service.RegisterCitizenInput
	if !bindRequest(c, &registerInput) {
		return
	}

	citizen, err := ah.AuthService.RegisterCitizen(c.Request.Context(), registerInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, citizen)
}

func (ah *Handler) RegisterCouncillor(c *gin.Context) {
	var registerInput service.RegisterCouncillorInput
	if !bindRequest(c, &registerInput) {
		return
	}

	councillor, err := ah.AuthService.RegisterCouncillor(c.Request.Context(), registerInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, councillor)
}

func (ah *Handler) RegisterOfficeMember(c *gin.Context) {
	var registerInput service.RegisterOfficeMemberInput
	if !bindRequest(c, &registerInput) {
		return
	}

	officeMember, err := ah.AuthService.RegisterOfficeMember(c.Request.Context(), registerInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, officeMember)
}

func (ah *Handler) PreviewOfficeMemberInvitation(c *gin.Context) {
	invitation, err := ah.AuthService.PreviewOfficeMemberInvitation(c.Request.Context(), c.Query("token"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, invitation)
}

func (ah *Handler) Me(c *gin.Context) {
	userID, _ := c.Get("userId")
	user, err := ah.AuthService.CurrentUser(c.Request.Context(), userID.(uint))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ah *Handler) Logout(c *gin.Context) {
	if cookie, err := c.Request.Cookie(refreshTokenCookieKey); err == nil {
		if logoutErr := ah.AuthService.Logout(c.Request.Context(), cookie.Value); logoutErr != nil {
			c.Error(logoutErr)
			return
		}
	}
	c.SetSameSite(http.SameSiteLaxMode)
	secure := c.Request.TLS != nil
	c.SetCookie(accessTokenCookieKey, "", -1, "/", "", secure, true)
	c.SetCookie(refreshTokenCookieKey, "", -1, "/", "", secure, true)
	c.Status(http.StatusNoContent)
}

func (ah *Handler) setAuthCookies(c *gin.Context, session *service.LoginOutput) {
	now := time.Now()
	c.SetSameSite(http.SameSiteLaxMode)
	secure := c.Request.TLS != nil
	c.SetCookie(accessTokenCookieKey, session.AccessToken, int(session.AccessTokenExpiresIn.Sub(now).Seconds()), "/", "", secure, true)
	c.SetCookie(refreshTokenCookieKey, session.RefreshToken, int(session.RefreshTokenExpiresIn.Sub(now).Seconds()), "/", "", secure, true)
}
