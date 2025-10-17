package main

import (
	"net/http"

	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/middleware"
	"github.com/WoWBytePaladin/go-mall/config"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()

	// 有了AccessLog 后, 就不需要gin.Logger这个中间件啦
	// g.Use(gin.Logger(), middleware.StartTrace())
	g.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	g.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		c.JSON(http.StatusOK, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})

	g.GET("/logger-test", func(c *gin.Context) {
		logger.New(c).Info("logger test", "key", "keyName", "val", 2)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	g.POST("/access-log-test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	g.GET("/panic-log-test", func(c *gin.Context) {
		var a map[string]string
		a["k"] = "v"
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   a,
		})
	})

	g.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
