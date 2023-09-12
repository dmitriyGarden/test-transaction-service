package main

import (
	"context"
	"log"

	"github.com/dmitriyGarden/test-transaction-service/adapter/in/brokers"
	"github.com/dmitriyGarden/test-transaction-service/adapter/in/web"
	"github.com/dmitriyGarden/test-transaction-service/adapter/out/persistence/pgres"
	"github.com/dmitriyGarden/test-transaction-service/app/service"
	"github.com/dmitriyGarden/test-transaction-service/pkg/config"
	"github.com/dmitriyGarden/test-transaction-service/pkg/logger"
	"github.com/nats-io/nats.go"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config.New: %v", err)
	}
	l := logger.New()
	storage, err := pgres.New(cfg)
	if err != nil {
		l.Fatal("pgres.New: ", err)
	}

	trService, err := service.New(cfg, storage)
	if err != nil {
		l.Fatal("service.New: ", err)
	}

	conn, err := nats.Connect(cfg.NatsConnectionString())
	if err != nil {
		l.Fatal("nats.Connect: ", err)
	}
	defer func(conn *nats.Conn) {
		_ = conn.Drain()
		conn.Close()
	}(conn)
	itr, err := brokers.GetBrokerNatsAdapter(conn, trService, l)
	if err != nil {
		l.Fatal("GetInternalNatsAdapter: ", err)
	}
	webAdapter, err := web.GetWebGrpcAdapter(cfg, trService, l)
	if err != nil {
		l.Fatal("web.GetWebGrpcAdapter: ", err)
	}
	go func() {
		err = itr.Run(context.Background())
		if err != nil {
			l.Fatal("itr.Run: ", err)
		}
	}()
	err = webAdapter.Run(context.Background())
	if err != nil {
		l.Fatal("webAdapter.Run: %v", err)
	}
	l.Infoln("Finished")
}
