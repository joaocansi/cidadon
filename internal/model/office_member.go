package model

import (
	"cidadon/internal/domain/entity"

	"gorm.io/gorm"
)

type OfficeMember struct {
	UserID   uint    `gorm:"primaryKey"`
	User     *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	OfficeID uint    `gorm:"not null;index"`
	Office   *Office `gorm:"foreignKey:OfficeID;constraint:OnDelete:CASCADE"`
	ImageURL string
}

func (om *OfficeMember) ToDomain() *entity.OfficeMember {
	if om == nil {
		return nil
	}

	return &entity.OfficeMember{
		UserID:   om.UserID,
		User:     om.User.ToDomain(),
		OfficeID: om.OfficeID,
		Office:   om.Office.ToDomain(),
		ImageURL: om.ImageURL,
	}
}

type OfficeMemberRequest struct {
	gorm.Model
	OfficeID uint    `gorm:"not null;index"`
	Office   *Office `gorm:"foreignKey:OfficeID;constraint:OnDelete:CASCADE"`
	Token    string  `gorm:"uniqueIndex;not null"`
	Email    string  `gorm:"not null"`
}

func (omr *OfficeMemberRequest) ToDomain() *entity.OfficeMemberRequest {
	if omr == nil {
		return nil
	}

	return &entity.OfficeMemberRequest{
		ID:        omr.ID,
		OfficeID:  omr.OfficeID,
		Office:    omr.Office.ToDomain(),
		Token:     omr.Token,
		Email:     omr.Email,
		CreatedAt: omr.CreatedAt,
		UpdatedAt: omr.UpdatedAt,
	}
}
