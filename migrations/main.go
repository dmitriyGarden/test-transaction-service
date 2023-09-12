package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/dmitriyGarden/test-transaction-service/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Expected at least 2 arguments")
	}
	command := os.Args[1]
	args := os.Args[2:]

	cfg, err := config.New()
	if err != nil {
		log.Fatal("config.New: ", err)
	}
	var coreDB *sql.DB
	if command != "create" {
		conn, err := sqlx.Connect("postgres", cfg.PostgresConnection())
		if err != nil {
			log.Fatal("sqlx.Connect: ", err)
		}

		coreDB = conn.DB
	}
	err = os.MkdirAll(cfg.MigrationsPath(), os.ModePerm)
	if err != nil {
		log.Fatalf("os.MkdirAll: %v", err)
	}
	goose.SetTableName("transaction_migrations")
	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatalf("goose.SetDialect: %v", err)
	}
	err = goose.Run(command, coreDB, cfg.MigrationsPath(), args...)
	if err != nil {
		log.Fatalf("goose.Run: %v", err)
	}
}
