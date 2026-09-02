package service

import (
	"cidadon/internal/domain/entity"
	"context"
)

type InviteOfficeMemberInput struct {
	Email string `json:"email" binding:"required,email"`
}
type InviteOfficeMemberOutput struct {
	Email     string `json:"email"`
	ExpiresAt string `json:"expires_at"`
}
type OfficeDirectoryItem struct {
	OfficeID       uint   `json:"office_id"`
	CouncillorName string `json:"councillor_name"`
	Party          string `json:"party"`
}
type PublicOfficeOutput struct {
	OfficeID       uint                         `json:"office_id"`
	Slug           string                       `json:"slug"`
	CouncillorName string                       `json:"councillor_name"`
	Party          string                       `json:"party"`
	ImageURL       string                       `json:"image_url"`
	City           string                       `json:"city"`
	State          string                       `json:"state"`
	Description    string                       `json:"description"`
	Contacts       []entity.OfficeContact       `json:"contacts"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks"`
}
type OfficeMemberOutput struct {
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageURL string `json:"image_url"`
}
type ManagedOfficeOutput struct {
	PublicOfficeOutput
	Members []OfficeMemberOutput `json:"members"`
}
type OfficeMemberRequestOutput struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

type CreateOfficeInput struct {
	CouncillorID   uint
	Description    string                       `json:"description" binding:"max=1500"`
	Contacts       []entity.OfficeContact       `json:"contacts" binding:"max=10,dive"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks" binding:"max=10,dive"`
}

type CreateOfficeOutput struct {
	OfficeID       uint                         `json:"office_id"`
	CouncillorID   uint                         `json:"councillor_id"`
	Description    string                       `json:"description"`
	Contacts       []entity.OfficeContact       `json:"contacts" binding:"max=10,dive"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks" binding:"max=10,dive"`
}

type UpdateOfficeInput struct {
	CouncillorID   uint
	Party          string                       `json:"party" binding:"required,min=2,max=80"`
	Description    string                       `json:"description" binding:"max=1500"`
	City           string                       `json:"city" binding:"required,min=2,max=120"`
	State          string                       `json:"state" binding:"required,len=2"`
	Contacts       []entity.OfficeContact       `json:"contacts" binding:"max=10,dive"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks" binding:"max=10,dive"`
}

type UpdateOfficeOutput struct {
	OfficeID       uint                         `json:"office_id"`
	CouncillorID   uint                         `json:"councillor_id"`
	Slug           string                       `json:"slug"`
	Party          string                       `json:"party"`
	Description    string                       `json:"description"`
	City           string                       `json:"city"`
	State          string                       `json:"state"`
	Contacts       []entity.OfficeContact       `json:"contacts"`
	SocialNetworks []entity.OfficeSocialNetwork `json:"social_networks"`
}

type OfficeService interface {
	Create(context.Context, CreateOfficeInput) (*CreateOfficeOutput, error)
	Update(context.Context, UpdateOfficeInput) (*UpdateOfficeOutput, error)
	InviteMember(context.Context, uint, InviteOfficeMemberInput) (*InviteOfficeMemberOutput, error)
	ResolveOfficeID(context.Context, uint, entity.UserRole) (uint, error)
	ListDirectory(context.Context, string, string) ([]OfficeDirectoryItem, error)
	SearchPublic(context.Context, string, string, string) ([]PublicOfficeOutput, error)
	FindPublic(context.Context, string) (*PublicOfficeOutput, error)
	FindManaged(context.Context, uint, entity.UserRole) (*ManagedOfficeOutput, error)
	ListMemberRequests(context.Context, uint) ([]OfficeMemberRequestOutput, error)
	CancelMemberRequest(context.Context, uint, uint) error
	RemoveMember(context.Context, uint, uint) error
}
