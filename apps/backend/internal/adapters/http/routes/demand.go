package routes

import (
	"cidadon/internal/adapters/http/handler"
	http "cidadon/internal/adapters/http/middleware"
	"cidadon/internal/domain/entity"

	"github.com/gin-gonic/gin"
)

func RegisterDemands(router *gin.Engine, middleware *http.AuthMiddleware, handler *handler.DemandHandler) {
	group := router.Group("/demands")
	group.GET("", handler.List)
	group.GET("/:id", handler.FindByID)
	group.GET("/:id/activity", middleware.AuthHandler(entity.CitizenUser, entity.CouncillorUser, entity.OfficeMemberUser), handler.Activity)
	group.POST("/:id/comments", middleware.AuthHandler(entity.CitizenUser, entity.CouncillorUser, entity.OfficeMemberUser), handler.Comment)
	citizen := group.Group("")
	citizen.Use(middleware.AuthHandler(entity.CitizenUser))
	citizen.POST("", handler.Create)
	citizen.GET("/mine", handler.ListMine)
	citizen.POST("/:id/confirm", handler.Confirm)
	citizen.POST("/:id/reopen", handler.Reopen)
	citizen.GET("/:id/support", handler.Support)
	citizen.PUT("/:id/support", handler.Support)
	citizen.DELETE("/:id/support", handler.Support)
	staff := group.Group("")
	staff.Use(middleware.AuthHandler(entity.CouncillorUser, entity.OfficeMemberUser))
	staff.GET("/office", handler.ListForOffice)
	staff.POST("/:id/claim", handler.Claim)
	staff.POST("/:id/start", handler.Start)
	staff.POST("/:id/request-confirmation", handler.RequestConfirmation)
	staff.POST("/:id/milestones", handler.CreateMilestone)
}
