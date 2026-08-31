package routes

import (
	"cidadon/internal/adapters/http/handler"
	http "cidadon/internal/adapters/http/middleware"
	"cidadon/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

func RegisterNotifications(router *gin.Engine, middleware *http.AuthMiddleware, handler *handler.DemandHandler) {
	group := router.Group("/notifications", middleware.AuthHandler(entity.CitizenUser, entity.CouncillorUser, entity.OfficeMemberUser))
	group.GET("", handler.Notifications)
	group.POST("/read", handler.ReadNotifications)
	group.GET("/stream", handler.NotificationStream)
	comments := router.Group("/comments")
	comments.POST("/:commentId/report", middleware.AuthHandler(entity.CitizenUser, entity.CouncillorUser, entity.OfficeMemberUser), handler.ReportComment)
	comments.POST("/:commentId/hide", middleware.AuthHandler(entity.CouncillorUser, entity.OfficeMemberUser), handler.HideComment)
}
