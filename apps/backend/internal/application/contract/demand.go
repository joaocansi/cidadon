package service

import (
	"cidadon/internal/domain/entity"
	"context"
	"time"
)

type CreateDemandInput struct {
	CitizenID        uint
	Title            string   `json:"title" form:"title" binding:"required,min=5,max=120"`
	Description      string   `json:"description" form:"description" binding:"required,min=10,max=1000"`
	Category         string   `json:"category" form:"category" binding:"required,min=3,max=80"`
	Street           string   `json:"street" form:"street" binding:"required,min=3,max=120"`
	Number           string   `json:"number" form:"number" binding:"omitempty,max=20"`
	Neighborhood     string   `json:"neighborhood" form:"neighborhood" binding:"required,min=2,max=120"`
	City             string   `json:"city" form:"city" binding:"required,min=2,max=120"`
	State            string   `json:"state" form:"state" binding:"required,min=2,max=2"`
	Latitude         float64  `json:"latitude" form:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude        float64  `json:"longitude" form:"longitude" binding:"required,gte=-180,lte=180"`
	Images           []string `json:"images" binding:"max=5"`
	DirectedOfficeID *uint    `json:"directed_office_id" form:"directed_office_id"`
}

type DemandOutput struct {
	ID                  uint                `json:"id"`
	Protocol            string              `json:"protocol"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	Category            string              `json:"category"`
	Street              string              `json:"street"`
	Number              string              `json:"number"`
	Neighborhood        string              `json:"neighborhood"`
	City                string              `json:"city"`
	State               string              `json:"state"`
	Status              entity.DemandStatus `json:"status"`
	SupportCount        int                 `json:"support_count"`
	CommentCount        int                 `json:"comment_count"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
	Latitude            float64             `json:"latitude"`
	Longitude           float64             `json:"longitude"`
	Images              []string            `json:"images"`
	DirectedOfficeID    *uint               `json:"directed_office_id,omitempty"`
	ResponsibleOfficeID *uint               `json:"responsible_office_id,omitempty"`
	ConfirmationDueAt   *time.Time          `json:"confirmation_due_at,omitempty"`
}
type DemandCommentInput struct {
	Body     string   `json:"body" form:"body" binding:"max=2000"`
	Images   []string `json:"images" binding:"max=5"`
	ParentID *uint    `json:"parent_id,omitempty" form:"parent_id"`
}

// DemandTimelineInput is deliberately separate from a comment. Status changes
// and milestones are part of the public operational record and therefore
// always require a written explanation.
type DemandTimelineInput struct {
	Message string   `json:"message" form:"message" binding:"required,min=3,max=2000"`
	Images  []string `json:"images" binding:"max=5"`
}

type DemandSupportOutput struct {
	SupportCount int  `json:"support_count"`
	Supported    bool `json:"supported"`
	CanSupport   bool `json:"can_support"`
}
type DemandActivityOutput struct {
	Events   []entity.DemandEvent   `json:"events"`
	Comments []entity.DemandComment `json:"comments"`
}
type NotificationOutput struct {
	ID        uint       `json:"id"`
	Type      string     `json:"type"`
	DemandID  uint       `json:"demand_id"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type DemandListFilters struct {
	City         string
	State        string
	Neighborhood string
	Category     string
	Status       entity.DemandStatus
}

type DemandService interface {
	Create(context.Context, CreateDemandInput) (*DemandOutput, error)
	FindByID(context.Context, uint) (*DemandOutput, error)
	List(context.Context, DemandListFilters) ([]DemandOutput, error)
	ListMine(context.Context, uint) ([]DemandOutput, error)
	ListForOffice(context.Context, uint, DemandListFilters) ([]DemandOutput, error)
	Claim(context.Context, uint, uint, uint, DemandTimelineInput) (*DemandOutput, error)
	Start(context.Context, uint, uint, uint, DemandTimelineInput) (*DemandOutput, error)
	RequestConfirmation(context.Context, uint, uint, uint, DemandTimelineInput) (*DemandOutput, error)
	Confirm(context.Context, uint, uint) (*DemandOutput, error)
	Reopen(context.Context, uint, uint, DemandTimelineInput) (*DemandOutput, error)
	CreateMilestone(context.Context, uint, uint, uint, DemandTimelineInput) error
	GetSupport(context.Context, uint, uint) (*DemandSupportOutput, error)
	AddSupport(context.Context, uint, uint) (*DemandSupportOutput, error)
	RemoveSupport(context.Context, uint, uint) (*DemandSupportOutput, error)
	Comment(context.Context, uint, uint, entity.UserRole, DemandCommentInput) (*entity.DemandComment, error)
	Activity(context.Context, uint) (*DemandActivityOutput, error)
	DeleteComment(context.Context, uint, uint) error
	ReportComment(context.Context, uint, uint, string) error
	HideComment(context.Context, uint, uint, uint) error
	ListNotifications(context.Context, uint) ([]NotificationOutput, error)
	ListNotificationsAfter(context.Context, uint, uint) ([]NotificationOutput, error)
	ReadNotifications(context.Context, uint, []uint) error
}
