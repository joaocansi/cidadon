package handler

import (
	"cidadon/internal/adapters/external/provider"
	service "cidadon/internal/application/contract"
	"cidadon/internal/application/usecase"
	"cidadon/internal/domain/entity"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DemandHandler struct {
	DemandService service.DemandService
	OfficeService service.OfficeService
	MediaService  *usecase.MediaService
}

func NewDemandHandler(demandService service.DemandService, officeService service.OfficeService, mediaService *usecase.MediaService) *DemandHandler {
	return &DemandHandler{
		DemandService: demandService,
		OfficeService: officeService,
		MediaService:  mediaService,
	}
}

func (h *DemandHandler) ListForOffice(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}
	role := c.MustGet("userRole").(entity.UserRole)
	officeID, err := h.OfficeService.ResolveOfficeID(c.Request.Context(), userID.(uint), role)
	if err != nil {
		c.Error(err)
		return
	}
	demands, err := h.DemandService.ListForOffice(c.Request.Context(), officeID, service.DemandListFilters{Status: entity.DemandStatus(c.Query("status"))})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, demands)
}

func (h *DemandHandler) Create(c *gin.Context) {
	var input service.CreateDemandInput
	if !bindMultipartRequest(c, &input) {
		return
	}

	citizenID, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}
	input.CitizenID = citizenID.(uint)
	stored, ok := h.storeImages(c, "demands")
	if !ok {
		return
	}
	input.Images = usecase.MediaURLs(stored)

	demand, err := h.DemandService.Create(c.Request.Context(), input)
	if err != nil {
		h.MediaService.DeleteAll(c.Request.Context(), stored)
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, demand)
}

func (h *DemandHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid demand id"))
		return
	}

	demand, err := h.DemandService.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, demand)
}

func (h *DemandHandler) List(c *gin.Context) {
	demands, err := h.DemandService.List(c.Request.Context(), service.DemandListFilters{
		City:         c.Query("city"),
		State:        c.Query("state"),
		Neighborhood: c.Query("neighborhood"),
		Category:     c.Query("category"),
		Status:       entity.DemandStatus(c.Query("status")),
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, demands)
}

func (h *DemandHandler) ListMine(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		c.Error(service.Unauthorized("not authorized"))
		return
	}
	demands, err := h.DemandService.ListMine(c.Request.Context(), userID.(uint))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, demands)
}

func (h *DemandHandler) Claim(c *gin.Context)               { h.officeAction(c, "claim") }
func (h *DemandHandler) Start(c *gin.Context)               { h.officeAction(c, "start") }
func (h *DemandHandler) RequestConfirmation(c *gin.Context) { h.officeAction(c, "confirm_request") }
func (h *DemandHandler) officeAction(c *gin.Context, action string) {
	id, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil {
		c.Error(service.InvalidInput("invalid demand id"))
		return
	}
	uid := c.MustGet("userId").(uint)
	role := c.MustGet("userRole").(entity.UserRole)
	office, e := h.OfficeService.ResolveOfficeID(c.Request.Context(), uid, role)
	if e != nil {
		c.Error(e)
		return
	}
	var input service.DemandTimelineInput
	if !bindMultipartRequest(c, &input) {
		return
	}
	stored, ok := h.storeImages(c, "timeline")
	if !ok {
		return
	}
	input.Images = usecase.MediaURLs(stored)
	var out *service.DemandOutput
	if action == "claim" {
		out, e = h.DemandService.Claim(c.Request.Context(), uint(id), office, uid, input)
	} else if action == "start" {
		out, e = h.DemandService.Start(c.Request.Context(), uint(id), office, uid, input)
	} else {
		out, e = h.DemandService.RequestConfirmation(c.Request.Context(), uint(id), office, uid, input)
	}
	if e != nil {
		h.MediaService.DeleteAll(c.Request.Context(), stored)
		c.Error(e)
		return
	}
	c.JSON(http.StatusOK, out)
}
func (h *DemandHandler) Confirm(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	out, e := h.DemandService.Confirm(c, uint(id), c.MustGet("userId").(uint))
	if e != nil {
		c.Error(e)
		return
	}
	c.JSON(http.StatusOK, out)
}
func (h *DemandHandler) Reopen(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in service.DemandTimelineInput
	if !bindMultipartRequest(c, &in) {
		return
	}
	stored, ok := h.storeImages(c, "timeline")
	if !ok {
		return
	}
	in.Images = usecase.MediaURLs(stored)
	out, e := h.DemandService.Reopen(c, uint(id), c.MustGet("userId").(uint), in)
	if e != nil {
		h.MediaService.DeleteAll(c.Request.Context(), stored)
		c.Error(e)
		return
	}
	c.JSON(http.StatusOK, out)
}
func (h *DemandHandler) CreateMilestone(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid demand id"))
		return
	}
	var input service.DemandTimelineInput
	if !bindMultipartRequest(c, &input) {
		return
	}
	stored, ok := h.storeImages(c, "timeline")
	if !ok {
		return
	}
	input.Images = usecase.MediaURLs(stored)
	userID := c.MustGet("userId").(uint)
	role := c.MustGet("userRole").(entity.UserRole)
	officeID, err := h.OfficeService.ResolveOfficeID(c.Request.Context(), userID, role)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.DemandService.CreateMilestone(c.Request.Context(), uint(id), officeID, userID, input); err != nil {
		h.MediaService.DeleteAll(c.Request.Context(), stored)
		c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}
func (h *DemandHandler) Support(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid demand id"))
		return
	}
	citizenID := c.MustGet("userId").(uint)
	var out *service.DemandSupportOutput
	switch c.Request.Method {
	case http.MethodGet:
		out, err = h.DemandService.GetSupport(c.Request.Context(), uint(id), citizenID)
	case http.MethodPut:
		out, err = h.DemandService.AddSupport(c.Request.Context(), uint(id), citizenID)
	case http.MethodDelete:
		out, err = h.DemandService.RemoveSupport(c.Request.Context(), uint(id), citizenID)
	default:
		c.Error(service.InvalidInput("unsupported support action"))
		return
	}
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, out)
}
func (h *DemandHandler) Comment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in service.DemandCommentInput
	if !bindMultipartRequest(c, &in) {
		return
	}
	stored, ok := h.storeImages(c, "comments")
	if !ok {
		return
	}
	in.Images = usecase.MediaURLs(stored)
	out, e := h.DemandService.Comment(c, uint(id), c.MustGet("userId").(uint), c.MustGet("userRole").(entity.UserRole), in)
	if e != nil {
		h.MediaService.DeleteAll(c.Request.Context(), stored)
		c.Error(e)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *DemandHandler) storeImages(c *gin.Context, prefix string) ([]provider.StoredMedia, bool) {
	stored, err := h.MediaService.StoreFiles(c.Request.Context(), prefix, multipartFiles(c, "images"), 5)
	if err != nil {
		c.Error(err)
		return nil, false
	}
	return stored, true
}
func (h *DemandHandler) Activity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid demand id"))
		return
	}
	out, e := h.DemandService.Activity(c, uint(id))
	if e != nil {
		c.Error(e)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *DemandHandler) ReportComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid comment id"))
		return
	}
	var input struct {
		Reason string `json:"reason" binding:"required,max=500"`
	}
	if !bindRequest(c, &input) {
		return
	}
	if err := h.DemandService.ReportComment(c.Request.Context(), uint(id), c.MustGet("userId").(uint), input.Reason); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DemandHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid comment id"))
		return
	}
	if err := h.DemandService.DeleteComment(c.Request.Context(), uint(id), c.MustGet("userId").(uint)); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DemandHandler) HideComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		c.Error(service.InvalidInput("invalid comment id"))
		return
	}
	userID := c.MustGet("userId").(uint)
	role := c.MustGet("userRole").(entity.UserRole)
	officeID, err := h.OfficeService.ResolveOfficeID(c.Request.Context(), userID, role)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.DemandService.HideComment(c.Request.Context(), uint(id), officeID, userID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *DemandHandler) Notifications(c *gin.Context) {
	out, err := h.DemandService.ListNotifications(c.Request.Context(), c.MustGet("userId").(uint))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, out)
}
func (h *DemandHandler) ReadNotifications(c *gin.Context) {
	var input struct {
		IDs []uint `json:"ids"`
	}
	if !bindRequest(c, &input) {
		return
	}
	if err := h.DemandService.ReadNotifications(c.Request.Context(), c.MustGet("userId").(uint), input.IDs); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *DemandHandler) NotificationStream(c *gin.Context) {
	userID := c.MustGet("userId").(uint)
	after := uint(0)
	if rawAfter := c.Query("after"); rawAfter != "" {
		parsed, err := strconv.ParseUint(rawAfter, 10, 64)
		if err != nil {
			c.Error(service.InvalidInput("invalid notifications cursor"))
			return
		}
		after = uint(parsed)
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache, no-transform")
	c.Header("X-Accel-Buffering", "no")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}
	_, _ = c.Writer.WriteString(": connected\n\n")
	flusher.Flush()
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		rows, err := h.DemandService.ListNotificationsAfter(c.Request.Context(), userID, after)
		if err != nil {
			return
		}
		for _, row := range rows {
			c.SSEvent("notification", row)
			after = row.ID
			flusher.Flush()
		}
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
		}
	}
}
