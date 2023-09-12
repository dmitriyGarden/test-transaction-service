package nats_server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dmitriyGarden/test-transaction-service/model"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Server struct {
	conn    *nats.Conn
	service model.ITransaction
	l       model.ILogger
}

func New(conn *nats.Conn, tr model.ITransaction, logger model.ILogger) (*Server, error) {
	return &Server{
		conn:    conn,
		service: tr,
		l:       logger,
	}, nil
}
func (s *Server) Run(ctx context.Context) error {
	sb, err := s.conn.Subscribe(servicePrefix+".*.*", s.handler)
	if err != nil {
		return fmt.Errorf("conn.Subscribe: %w", err)
	}
	select {
	case <-ctx.Done():
		_ = sb.Unsubscribe()
	}
	return nil
}

func (s *Server) handler(msg *nats.Msg) {
	topic := new(subject)
	err := topic.parse(msg.Subject)
	resp := respMessage{
		Type: successMessage,
	}
	if err != nil {
		resp.Type = errorMessage
		resp.Payload = err.Error()
		err = s.sendMessage(msg, resp)
		if err != nil {
			s.l.Error("sendMessage: ", err)
		}
		return
	}
	switch topic.method {
	case transactionBalanceMethod:
		res, err := s.GetBalance(msg.Data)
		if err != nil {
			resp.Type = errorMessage
			resp.Payload = err.Error()
		} else {
			resp.Payload = res
		}
	default:
		resp.Type = errorMessage
		resp.Payload = fmt.Sprintf("Unexpected method: %s", topic.method)
	}
	s.l.Debug("Response: ", resp)
	err = s.sendMessage(msg, resp)
	if err != nil {
		s.l.Error("sendMessage: ", err)
	}
}

func (s *Server) sendMessage(msg *nats.Msg, data respMessage) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	err = msg.Respond(body)
	if err != nil {
		return fmt.Errorf("msg.Respond: %w", err)
	}
	return nil
}

func (s *Server) reqID(ctx context.Context) string {
	id, ok := ctx.Value("requestID").(string)
	if ok && id != "" {
		return id
	}
	return uuid.NewString()
}
