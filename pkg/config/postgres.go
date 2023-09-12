package config

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

func (c *Config) PostgresConnection() string {
	return fmt.Sprintf("postgresql://%s:%s/%s?user=%s&password=%s&sslmode=disable",
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_PORT"),
		url.PathEscape(viper.GetString("DATABASE_NAME")),
		url.QueryEscape(viper.GetString("DATABASE_USER")),
		url.QueryEscape(viper.GetString("DATABASE_PASSWORD")),
	)
}

func (c *Config) MigrationsPath() string {
	return viper.GetString("migration_path")
}
