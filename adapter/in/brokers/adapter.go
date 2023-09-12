package brokers

import (
	"github.com/dmitriyGarden/test-transaction-service/adapter/in/brokers/nats_server"
	"github.com/dmitriyGarden/test-transaction-service/model"
	"github.com/nats-io/nats.go"
)

func GetBrokerNatsAdapter(conn *nats.Conn, tr model.ITransaction, logger model.ILogger) (model.IWebAdapter, error) {
	return nats_server.New(conn, tr, logger)
}
