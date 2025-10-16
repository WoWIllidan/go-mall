package main

import (
	"net/http"

	"github.com/WoWBytePaladin/go-mall/config"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/config-read", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"type":     config.Database.Type,
			"max_life": config.Database.MaxLifeTime,
		})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
