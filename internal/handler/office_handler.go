package handler

import (
	"cidadon/internal/service"

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

}

func (ah *OfficeHandler) Update(c *gin.Context) {

}

func (ah *OfficeHandler) NewMemberRequest(c *gin.Context) {

}
