package service

import (
	"cidadon/internal/domain/entity"
	"context"
	"time"
)

type CreateDemandInput struct {
	CitizenID        uint
	Title            string   `json:"title" binding:"required,min=5,max=120"`
	Description      string   `json:"description" binding:"required,min=10,max=1000"`
	Category         string   `json:"category" binding:"required,min=3,max=80"`
	Street           string   `json:"street" binding:"required,min=3,max=120"`
	Number           string   `json:"number" binding:"omitempty,max=20"`
	Neighborhood     string   `json:"neighborhood" binding:"required,min=2,max=120"`
	City             string   `json:"city" binding:"required,min=2,max=120"`
	State            string   `json:"state" binding:"required,min=2,max=2"`
	Latitude         float64  `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude        float64  `json:"longitude" binding:"required,gte=-180,lte=180"`
	Images           []string `json:"images" binding:"max=5"`
	DirectedOfficeID *uint    `json:"directed_office_id"`
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
	Body   string   `json:"body" binding:"max=2000"`
	Images []string `json:"images" binding:"max=5"`
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

type UpdateDemandStatusInput struct {
	Status entity.DemandStatus `json:"status" binding:"required"`
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
	UpdateStatus(context.Context, uint, uint, UpdateDemandStatusInput) (*DemandOutput, error)
	Claim(context.Context, uint, uint, uint) (*DemandOutput, error)
	Start(context.Context, uint, uint, uint) (*DemandOutput, error)
	RequestConfirmation(context.Context, uint, uint, uint, DemandCommentInput) (*DemandOutput, error)
	Confirm(context.Context, uint, uint) (*DemandOutput, error)
	Reopen(context.Context, uint, uint, DemandCommentInput) (*DemandOutput, error)
	Comment(context.Context, uint, uint, entity.UserRole, DemandCommentInput) (*entity.DemandComment, error)
	Activity(context.Context, uint) (*DemandActivityOutput, error)
	ReportComment(context.Context, uint, uint, string) error
	HideComment(context.Context, uint, uint, uint) error
	ListNotifications(context.Context, uint) ([]NotificationOutput, error)
	ReadNotifications(context.Context, uint, []uint) error
}
