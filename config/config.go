package config

import "time"

var (
	App      *appConfig
	Database *databaseConfig
)

type appConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type databaseConfig struct {
	Type        string        `mapstructure:"type"`
	DSN         string        `mapstructure:"dsn"`
	MaxOpenConn int           `mapstructure:"maxopen"`
	MaxIdleConn int           `mapstructure:"maxidle"`
	MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
}
