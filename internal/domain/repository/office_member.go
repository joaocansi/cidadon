package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateOfficeMemberData struct {
	OfficeID uint
	Name     string
	Email    string
	Password string
	ImageURL string
}

type OfficeMemberRepository interface {
	Create(context.Context, CreateOfficeMemberData) (*entity.OfficeMember, error)
}

type CreateOfficeMemberRequestData struct {
	OfficeID uint
	Email    string
	Token    string
}

type OfficeMemberRequestRepository interface {
	Create(context.Context, CreateOfficeMemberRequestData) (*entity.OfficeMemberRequest, error)
	FindByToken(context.Context, string) (*entity.OfficeMemberRequest, error)
	Delete(context.Context, uint) error
}
