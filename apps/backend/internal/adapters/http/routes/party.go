package routes

import (
	"cidadon/internal/adapters/http/handler"

	"github.com/gin-gonic/gin"
)

func RegisterParties(router *gin.Engine, handler *handler.PartyHandler) {
	router.GET("/parties", handler.List)
}
