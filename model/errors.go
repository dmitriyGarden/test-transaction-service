package model

import "fmt"

var (
	ErrNotFound     = fmt.Errorf("not found")
	ErrAuthRequired = fmt.Errorf("auth required")
	ErrInvalidToken = fmt.Errorf("invalid token")
)
