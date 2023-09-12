package nats_server

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Server) GetBalance(payload []byte) (int64, error) {
	uid := uuid.New()
	err := uid.UnmarshalBinary(payload)
	if err != nil {
		return 0, fmt.Errorf("uid.UnmarshalBinary: %w", err)
	}
	s.l.Debug("GetBalance: ", uid.String())
	res, err := s.service.GetBalance(context.Background(), uid)
	if err != nil {
		return 0, fmt.Errorf("service.GetBalance: %w", err)
	}
	return res, nil
}
