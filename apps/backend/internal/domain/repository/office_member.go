package repository

import (
	"cidadon/internal/domain/entity"
	"context"
	"time"
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
	FindByUserID(context.Context, uint) (*entity.OfficeMember, error)
	ListByOfficeID(context.Context, uint) ([]entity.OfficeMember, error)
}

type CreateOfficeMemberRequestData struct {
	OfficeID  uint
	Email     string
	Token     string
	ExpiresAt time.Time
}

type OfficeMemberRequestRepository interface {
	Create(context.Context, CreateOfficeMemberRequestData) (*entity.OfficeMemberRequest, error)
	FindByToken(context.Context, string) (*entity.OfficeMemberRequest, error)
	FindByID(context.Context, uint) (*entity.OfficeMemberRequest, error)
	ListByOfficeID(context.Context, uint) ([]entity.OfficeMemberRequest, error)
	Delete(context.Context, uint) error
	FindByEmailAndOffice(context.Context, string, uint) (*entity.OfficeMemberRequest, error)
}
