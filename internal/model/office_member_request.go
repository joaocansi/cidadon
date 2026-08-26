package model

import (
	"cidadon/internal/domain/entity"

	"gorm.io/gorm"
)

type OfficeMember struct {
	gorm.Model
	OfficeID uint
	Office   *Office
	UserID   uint
	User     *User
	ImageURL string
}

func (om *OfficeMember) ToDomain() *entity.OfficeMember {
	if om == nil {
		return nil
	}

	return &entity.OfficeMember{
		ID:        om.ID,
		OfficeID:  om.OfficeID,
		Office:    om.Office.ToDomain(),
		UserID:    om.UserID,
		User:      om.User.ToDomain(),
		ImageURL:  om.ImageURL,
		CreatedAt: om.CreatedAt,
		UpdatedAt: om.UpdatedAt,
	}
}

type OfficeMemberRequest struct {
	gorm.Model
	OfficeID uint
	Office   *Office
	Token    string
	Email    string
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
