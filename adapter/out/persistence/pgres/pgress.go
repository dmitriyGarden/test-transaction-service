package pgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type IConfig interface {
	PostgresConnection() string
}

type DB struct {
	db *sqlx.DB
}

func (d *DB) init(cfg IConfig) error {
	var db *sqlx.DB
	var err error
	str := cfg.PostgresConnection()
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("postgres", str)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return fmt.Errorf("sqlx.Connect: %w", err)
	}
	// Setting timeout to flush out the dead connections
	// ref: https://stackoverflow.com/questions/50338338/intermittent-connection-reset-by-peer-sql-postgres
	db.SetConnMaxLifetime(time.Minute * 15)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Minute)
	db.SetMaxOpenConns(20)
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}
	d.db = db
	return nil
}

func New(cfg IConfig) (*DB, error) {
	db := new(DB)
	err := db.init(cfg)
	if err != nil {
		return nil, fmt.Errorf("db.init: %w", err)
	}
	return db, nil
}
