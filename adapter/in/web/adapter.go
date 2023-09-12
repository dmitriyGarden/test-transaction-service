package web

import (
	"github.com/dmitriyGarden/test-transaction-service/adapter/in/web/grpc_server/server"
	"github.com/dmitriyGarden/test-transaction-service/model"
)

func GetWebGrpcAdapter(cfg server.IConfig, tr model.ITransaction, l model.ILogger) (model.IWebAdapter, error) {
	return server.New(cfg, tr, l)
}
