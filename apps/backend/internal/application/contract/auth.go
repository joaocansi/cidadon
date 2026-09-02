package service

import (
	"cidadon/internal/domain/entity"
	"context"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, data LoginInput) (*LoginOutput, error)
	RegisterCitizen(ctx context.Context, data RegisterCitizenInput) (*RegisterCitizenOutput, error)
	RegisterCouncillor(ctx context.Context, data RegisterCouncillorInput) (*RegisterCouncillorOutput, error)
	RegisterOfficeMember(ctx context.Context, data RegisterOfficeMemberInput) (*RegisterOfficeMemberOutput, error)
	PreviewOfficeMemberInvitation(ctx context.Context, token string) (*OfficeMemberInvitationOutput, error)
	CurrentUser(ctx context.Context, userID uint) (*CurrentUserOutput, error)
	Logout(ctx context.Context, refreshToken string) error
}

type CurrentUserOutput struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Role     entity.UserRole `json:"role"`
	ImageURL string          `json:"image_url,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=72"`
}

type LoginOutput struct {
	AccessToken           string          `json:"access_token"`
	RefreshToken          string          `json:"refresh_token"`
	AccessTokenExpiresIn  time.Time       `json:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Time       `json:"refresh_token_expires_in"`
	Role                  entity.UserRole `json:"role"`
}
type RegisterBaseInput struct {
	Name     string `json:"name" form:"name" binding:"required,min=3,max=120"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,gte=6,lte=72"`
}

type RegisterBaseOutput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterCitizenInput struct {
	RegisterBaseInput
	City  string `json:"city" form:"city" binding:"required,min=2,max=120"`
	State string `json:"state" form:"state" binding:"required,len=2"`
}

type RegisterCitizenOutput struct {
	RegisterBaseOutput
	City  string `json:"city"`
	State string `json:"state"`
}
type RegisterCouncillorInput struct {
	RegisterBaseInput
	Party    string `json:"party" form:"party" binding:"required,min=2,max=80"`
	ImageURL string `json:"image_url" form:"-"`
	City     string `json:"city" form:"city" binding:"required,min=2,max=120"`
	State    string `json:"state" form:"state" binding:"required,len=2"`
}

type RegisterCouncillorOutput struct {
	RegisterBaseOutput
	ImageURL string `json:"image_url"`
	Party    string `json:"party"`
	City     string `json:"city"`
	State    string `json:"state"`
}
type RegisterOfficeMemberInput struct {
	Name     string `json:"name" form:"name" binding:"required,min=3,max=120"`
	Password string `json:"password" form:"password" binding:"required,gte=6,lte=72"`
	Token    string `json:"token" form:"token" binding:"required,min=32,max=256"`
	ImageURL string `json:"image_url" form:"-"`
}

type RegisterOfficeMemberOutput struct {
	RegisterBaseOutput
	OfficeID uint   `json:"office_id"`
	ImageURL string `json:"image_url"`
}

type OfficeMemberInvitationOutput struct {
	OfficeID       uint   `json:"office_id"`
	CouncillorName string `json:"councillor_name"`
	Party          string `json:"party"`
	ImageURL       string `json:"image_url"`
	City           string `json:"city"`
	State          string `json:"state"`
	ExpiresAt      string `json:"expires_at"`
}
