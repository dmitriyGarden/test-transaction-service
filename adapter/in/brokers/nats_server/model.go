package nats_server

import (
	"fmt"
	"strings"
)

const servicePrefix = "transaction"

type transactionMethod string

const (
	transactionBalanceMethod transactionMethod = "balance"
)

type messageType string

const (
	errorMessage   messageType = "error"
	successMessage messageType = "success"
)

type subject struct {
	requestID string
	method    transactionMethod
}

func (t *subject) String() string {
	return servicePrefix + "." + string(t.method) + "." + t.requestID
}

func (t *subject) parse(topic string) error {
	arr := strings.Split(topic, ".")
	if len(arr) != 3 {
		return fmt.Errorf("unexpected subject format. expected format: \"%s.method.requiestid\"", servicePrefix)
	}
	t.method = transactionMethod(arr[1])
	t.requestID = arr[2]
	return nil
}

type response struct {
	Type    messageType
	Payload interface{}
}
