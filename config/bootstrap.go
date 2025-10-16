package config

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

//go:embed *.yaml
var configs embed.FS

func init() {
	env, ok := os.LookupEnv("ENV")
	if !ok || env == "" {
		panic("env variable not set")
	}

	vp := viper.New()

	configFileStream, err := configs.ReadFile("application." + env + ".yaml")
	fmt.Println("configFileStream", string(configFileStream))

	if err != nil {
		panic(err)
	}

	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewReader(configFileStream))
	if err != nil {
		panic(err)
	}

	err = vp.UnmarshalKey("app", &App)
	if err != nil {
		panic(err)
	}

	err = vp.UnmarshalKey("database", &Database)
	if err != nil {
		panic(err)
	}

	Database.MaxLifeTime *= time.Second
}
