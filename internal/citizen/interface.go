package citizen

import "context"

type Repository interface {
	Create(context.Context, *Citizen) error
}
