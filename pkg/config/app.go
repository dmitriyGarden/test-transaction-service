package config

import "github.com/spf13/viper"

func (c *Config) GetListen() string {
	return viper.GetString("listen")
}

func (c *Config) JWTSecret() []byte {
	return []byte(viper.GetString("JWT_SECRET"))
}
