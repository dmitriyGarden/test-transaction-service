package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func (c *Config) NatsConnectionString() string {
	return fmt.Sprintf(
		"nats://%s:%s@%s:%s",
		viper.GetString("NATS_USER"),
		viper.GetString("NATS_PASSWORD"),
		viper.GetString("NATS_HOST"),
		viper.GetString("NATS_PORT"),
	)
}
