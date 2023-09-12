package model

import (
	"context"

	"github.com/google/uuid"
)

type ITransaction interface {
	GetBalance(ctx context.Context, uid uuid.UUID) (int64, error)
	AddBalance(ctx context.Context, uid uuid.UUID, amount int64) error
	GetUserFromToken(token string) (uuid.UUID, error)
}
