package database

import "context"

type dbContextKey struct{}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fc func(ctx context.Context) error) error
}
