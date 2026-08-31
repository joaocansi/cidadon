package handler

import (
	service "cidadon/internal/application/contract"
	"cidadon/internal/domain/entity"
	"net/http"
	"strconv"

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
	if !bindRequest(c, &createOfficeInput) {
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
	if !bindRequest(c, &updateOfficeInput) {
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
	var input service.InviteOfficeMemberInput
	if !bindRequest(c, &input) {
		return
	}
	userID, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}
	invite, err := ah.OfficeService.InviteMember(c.Request.Context(), userID.(uint), input)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, invite)
}

func (ah *OfficeHandler) ListMemberRequests(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}
	requests, err := ah.OfficeService.ListMemberRequests(c.Request.Context(), userID.(uint))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, requests)
}

func (ah *OfficeHandler) CancelMemberRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid invitation id"))
		return
	}
	userID, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}
	if err := ah.OfficeService.CancelMemberRequest(c.Request.Context(), userID.(uint), uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (ah *OfficeHandler) RemoveMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid member id"))
		return
	}
	userID, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}
	if err := ah.OfficeService.RemoveMember(c.Request.Context(), userID.(uint), uint(id)); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (ah *OfficeHandler) ListDirectory(c *gin.Context) {
	items, err := ah.OfficeService.ListDirectory(c.Request.Context(), c.Query("city"), c.Query("state"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ah *OfficeHandler) SearchPublic(c *gin.Context) {
	items, err := ah.OfficeService.SearchPublic(c.Request.Context(), c.Query("q"), c.Query("city"), c.Query("state"))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ah *OfficeHandler) FindPublic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid office id"))
		return
	}
	item, err := ah.OfficeService.FindPublic(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ah *OfficeHandler) FindManaged(c *gin.Context) {
	userID, _ := c.Get("userId")
	role, _ := c.Get("userRole")
	office, err := ah.OfficeService.FindManaged(c.Request.Context(), userID.(uint), role.(entity.UserRole))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, office)
}
