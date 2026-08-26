package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type OfficeMemberRepository interface {
	Create(context.Context, entity.OfficeMember) (entity.OfficeMember, error)
}
