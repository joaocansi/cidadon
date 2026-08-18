package councillor

import "context"

type Repository interface {
	Create(context.Context, *Councillor) error
}
