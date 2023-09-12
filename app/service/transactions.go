package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmitriyGarden/test-transaction-service/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type IConfig interface {
	JWTSecret() []byte
}

type TransactionService struct {
	cfg     IConfig
	storage model.IStorage
}

func New(cfg IConfig, s model.IStorage) (*TransactionService, error) {
	return &TransactionService{
		cfg:     cfg,
		storage: s,
	}, nil
}

func (c *TransactionService) GetBalance(ctx context.Context, uid uuid.UUID) (int64, error) {
	res, err := c.storage.GetBalance(ctx, uid)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return 0, fmt.Errorf("storage.GetBalance: %w", err)
	}
	return res, nil
}

func (c *TransactionService) AddBalance(ctx context.Context, uid uuid.UUID, amount int64) error {
	//TODO do not allow negative balance?
	return c.storage.AddBalance(ctx, uid, amount)
}

func (c *TransactionService) getJWTSecret(token *jwt.Token) (interface{}, error) {
	if token.Method != jwt.SigningMethodHS256 {
		return nil, fmt.Errorf("unexpected method %s. %w", token.Method.Alg(), model.ErrInvalidToken)
	}
	return c.cfg.JWTSecret(), nil
}

func (c *TransactionService) GetUserFromToken(token string) (uuid.UUID, error) {
	claims := new(jwtClaim)
	t, err := jwt.ParseWithClaims(token, claims, c.getJWTSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%v. %w", err, model.ErrInvalidToken)
	}
	if !t.Valid {
		return uuid.Nil, model.ErrInvalidToken
	}
	uid, err := uuid.Parse(claims.UID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%v. %w", err, model.ErrInvalidToken)
	}
	return uid, nil
}

type jwtClaim struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}
