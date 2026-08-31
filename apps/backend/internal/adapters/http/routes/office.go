package routes

import (
	"cidadon/internal/adapters/http/handler"
	http "cidadon/internal/adapters/http/middleware"
	"cidadon/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

func RegisterOffice(router *gin.Engine, middleware *http.AuthMiddleware, handler *handler.OfficeHandler) {
	group := router.Group("/office")
	group.GET("", handler.SearchPublic)
	group.GET("/me", middleware.AuthHandler(entity.CouncillorUser, entity.OfficeMemberUser), handler.FindManaged)
	group.POST("", middleware.AuthHandler(entity.CouncillorUser), handler.Create)
	group.POST("/member-request", middleware.AuthHandler(entity.CouncillorUser), handler.NewMemberRequest)
	group.GET("/member-requests", middleware.AuthHandler(entity.CouncillorUser), handler.ListMemberRequests)
	group.DELETE("/member-requests/:id", middleware.AuthHandler(entity.CouncillorUser), handler.CancelMemberRequest)
	group.DELETE("/members/:id", middleware.AuthHandler(entity.CouncillorUser), handler.RemoveMember)
	group.PUT("", middleware.AuthHandler(entity.CouncillorUser), handler.Update)
	group.GET("/:id", handler.FindPublic)
}
