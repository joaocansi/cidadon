package service

import (
	"context"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, data LoginInput) (*LoginOutput, error)
	RegisterCitizen(ctx context.Context, data RegisterCitizenInput) (*RegisterCitizenOutput, error)
	RegisterCouncillor(ctx context.Context, data RegisterCouncillorInput) (*RegisterCouncillorOutput, error)
	RegisterOfficeMember(ctx context.Context, data RegisterOfficeMemberInput) (*RegisterOfficeMemberOutput, error)
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=12"`
}

type LoginOutput struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresIn  time.Time `json:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Time `json:"refresh_token_expires_in"`
}
type RegisterBaseInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=12"`
}

type RegisterBaseOutput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterCitizenInput struct {
	RegisterBaseInput
	City  string `json:"city" binding:"required"`
	State string `json:"state" binding:"required"`
}

type RegisterCitizenOutput struct {
	RegisterBaseOutput
	City  string `json:"city"`
	State string `json:"state"`
}
type RegisterCouncillorInput struct {
	RegisterBaseInput
	Party    string `json:"party" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
	City     string `json:"city" binding:"required"`
	State    string `json:"state" binding:"required"`
}

type RegisterCouncillorOutput struct {
	RegisterBaseOutput
	ImageURL string `json:"image_url"`
	Party    string `json:"party"`
	City     string `json:"city"`
	State    string `json:"state"`
}
type RegisterOfficeMemberInput struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=12"`
	Token    string `json:"token" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
}

type RegisterOfficeMemberOutput struct {
	RegisterBaseOutput
	OfficeID uint   `json:"office_id"`
	ImageURL string `json:"image_url"`
}
