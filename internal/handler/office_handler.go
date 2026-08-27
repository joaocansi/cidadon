package handler

import (
	"cidadon/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OfficeHandler struct {
	OfficeService service.OfficeService
}

func NewOfficeHandler(officeService service.OfficeService) *OfficeHandler {
	return &OfficeHandler{
		OfficeService: officeService,
	}
}

func (ah *OfficeHandler) Create(c *gin.Context) {
	var createOfficeInput service.CreateOfficeInput
	if err := c.ShouldBind(&createOfficeInput); err != nil {
		c.Error(err)
		return
	}

	councillorId, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}

	createOfficeInput.CouncillorID = councillorId.(uint)
	createdOffice, err := ah.OfficeService.Create(c.Request.Context(), createOfficeInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, createdOffice)
}

func (ah *OfficeHandler) Update(c *gin.Context) {
	var updateOfficeInput service.UpdateOfficeInput
	if err := c.ShouldBind(&updateOfficeInput); err != nil {
		c.Error(err)
		return
	}

	councillorId, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}

	updateOfficeInput.CouncillorID = councillorId.(uint)
	updatedOffice, err := ah.OfficeService.Update(c.Request.Context(), updateOfficeInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedOffice)
}

func (ah *OfficeHandler) NewMemberRequest(c *gin.Context) {

}
