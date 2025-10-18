package main

import (
	"github.com/WoWBytePaladin/go-mall/api/router"
	"github.com/WoWBytePaladin/go-mall/common/enum"
	"github.com/WoWBytePaladin/go-mall/config"
	"github.com/gin-gonic/gin"
)

func main() {
	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()

	router.RegisterRoutes(g)

	g.Run(":8080")

}
