package routes

import (
	"cidadon/internal/adapters/http/handler"
	http "cidadon/internal/adapters/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuth(router *gin.Engine, middleware *http.AuthMiddleware, handler *handler.Handler) {
	group := router.Group("/auth")
	group.POST("/login", handler.Login)
	group.POST("/register/citizen", handler.RegisterCitizen)
	group.POST("/register/councillor", handler.RegisterCouncillor)
	group.GET("/register/office-member/invitation", handler.PreviewOfficeMemberInvitation)
	group.POST("/register/office-member", handler.RegisterOfficeMember)
	group.GET("/me", middleware.OptionalAuthHandler(), handler.Me)
	group.POST("/logout", handler.Logout)
}
