package config

import (
	"bytes"
	"embed"
	"github.com/spf13/viper"
	"os"
	"time"
)

// **嵌入文件只能在写embed指令的Go文件的同级目录或者子目录中
//go:embed *.yaml
var configs embed.FS

func init() {
	env := os.Getenv("ENV")
	vp := viper.New()
	// 根据环境变量 ENV 决定要读取的应用启动配置
	configFileStream, err := configs.ReadFile("application." + env + ".yaml")
	if err != nil {
		panic(err)
	}
	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewReader(configFileStream))
	if err != nil {
		// 加载不到应用配置, 阻挡应用的继续启动
		panic(err)
	}
	vp.UnmarshalKey("app", &App)

	vp.UnmarshalKey("database", &Database)
	Database.MaxLifeTime *= time.Second
}
