package model

import (
	"context"

	"github.com/google/uuid"
)

type IStorage interface {
	GetBalance(ctx context.Context, uid uuid.UUID) (int64, error)
	AddBalance(ctx context.Context, uid uuid.UUID, sum int64) error
}
