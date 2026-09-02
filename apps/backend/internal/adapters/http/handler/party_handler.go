package handler

import (
	"cidadon/internal/domain/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PartyHandler struct{}

func NewPartyHandler() *PartyHandler { return &PartyHandler{} }

func (h *PartyHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, entity.OfficialParties())
}
