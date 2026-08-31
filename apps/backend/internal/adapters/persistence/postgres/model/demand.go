package model

import (
	"cidadon/internal/domain/entity"
	"encoding/json"

	"gorm.io/gorm"
	"time"
)

type Demand struct {
	gorm.Model
	Protocol            string   `gorm:"uniqueIndex;not null"`
	CitizenID           uint     `gorm:"not null;index"`
	Citizen             *Citizen `gorm:"foreignKey:CitizenID;references:UserID;constraint:OnDelete:CASCADE"`
	Title               string   `gorm:"not null"`
	Description         string   `gorm:"not null"`
	Category            string   `gorm:"not null;index"`
	Street              string   `gorm:"not null"`
	Number              string
	Neighborhood        string          `gorm:"not null;index"`
	City                string          `gorm:"not null;index"`
	State               string          `gorm:"not null;index"`
	Latitude            float64         `gorm:"not null;default:0;index"`
	Longitude           float64         `gorm:"not null;default:0;index"`
	Images              json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	DirectedOfficeID    *uint           `gorm:"index"`
	ResponsibleOfficeID *uint           `gorm:"index"`
	ClaimedByUserID     *uint           `gorm:"index"`
	ConfirmationDueAt   *time.Time      `gorm:"index"`
	Status              string          `gorm:"not null;default:registered;index"`
	SupportCount        int             `gorm:"not null;default:0"`
	CommentCount        int             `gorm:"not null;default:0"`
}

type DemandEvent struct {
	gorm.Model
	DemandID    uint            `gorm:"not null;index"`
	Type        string          `gorm:"not null;index"`
	ActorUserID *uint           `gorm:"index"`
	Metadata    json.RawMessage `gorm:"type:jsonb;not null;default:'{}'"`
}
type DemandComment struct {
	gorm.Model
	DemandID       uint            `gorm:"not null;index"`
	AuthorID       uint            `gorm:"not null;index"`
	Body           string          `gorm:"type:text"`
	Images         json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	HiddenAt       *time.Time      `gorm:"index"`
	HiddenByUserID *uint
}
type DemandCommentReport struct {
	gorm.Model
	CommentID  uint   `gorm:"not null;uniqueIndex:idx_comment_reporter"`
	ReporterID uint   `gorm:"not null;uniqueIndex:idx_comment_reporter"`
	Reason     string `gorm:"type:text"`
}
type Notification struct {
	gorm.Model
	UserID   uint   `gorm:"not null;index"`
	Type     string `gorm:"not null"`
	DemandID uint   `gorm:"not null;index"`
	EventID  *uint
	ReadAt   *time.Time `gorm:"index"`
}

type DemandAssignment struct {
	gorm.Model
	DemandID uint `gorm:"not null;uniqueIndex:idx_demand_office"`
	OfficeID uint `gorm:"not null;uniqueIndex:idx_demand_office;index"`
}

func (d *Demand) ToDomain() *entity.Demand {
	if d == nil {
		return nil
	}

	images := make([]string, 0)
	_ = json.Unmarshal(d.Images, &images)
	return &entity.Demand{
		ID:           d.ID,
		Protocol:     d.Protocol,
		CitizenID:    d.CitizenID,
		Citizen:      demandCitizenToDomain(d.Citizen),
		Title:        d.Title,
		Description:  d.Description,
		Category:     d.Category,
		Street:       d.Street,
		Number:       d.Number,
		Neighborhood: d.Neighborhood,
		City:         d.City,
		State:        d.State,
		Latitude:     d.Latitude, Longitude: d.Longitude, Images: images, DirectedOfficeID: d.DirectedOfficeID,
		ResponsibleOfficeID: d.ResponsibleOfficeID, ClaimedByUserID: d.ClaimedByUserID, ConfirmationDueAt: d.ConfirmationDueAt,
		Status:       entity.DemandStatus(d.Status),
		SupportCount: d.SupportCount,
		CommentCount: d.CommentCount,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func demandCitizenToDomain(citizen *Citizen) *entity.Citizen {
	if citizen == nil {
		return nil
	}
	return citizen.ToDomain()
}
