package config

import "time"

// 项目通过这里的变量读取应用配置中的对应项
var (
	App      *appConfig
	Database *databaseConfig
	Redis    *redisConfig
)

// App配置
type appConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Log  struct {
		FilePath         string `mapstructure:"path"`
		FileMaxSize      int    `mapstructure:"max_size"`
		BackUpFileMaxAge int    `mapstructure:"max_age"`
	}
	Pagination struct {
		DefaultSize int `mapstructure:"default_size"`
		MaxSize     int `mapstructure:"max_size"`
	}
}

// 数据库配置
type databaseConfig struct {
	Master DbConnectOption `mapstructure:"master"`
	Slave  DbConnectOption `mapstructure:"slave"`
}

type DbConnectOption struct {
	Type        string        `mapstructure:"type"`
	DSN         string        `mapstructure:"dsn"`
	MaxOpenConn int           `mapstructure:"maxopen"`
	MaxIdleConn int           `mapstructure:"maxidle"`
	MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
}

// Redis 配置
type redisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
	DB       int    `mapstructure:"db"`
}
