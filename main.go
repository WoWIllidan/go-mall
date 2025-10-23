package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/WoWBytePaladin/go-mall/api/router"
	"github.com/WoWBytePaladin/go-mall/common/enum"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/config"
	"github.com/gin-gonic/gin"
)

func main() {
	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.New()
	router.RegisterRoutes(g)
	server := http.Server{
		Addr:    ":8080",
		Handler: g,
	}

	log := logger.New(context.Background())

	// 创建系统信号接收器
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("ShutdownServerError", "err", err)
		}
	}()

	log.Info("Starting GO MALL HTTP server...")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			// 服务正常收到关闭信号后Close
			log.Info("Server closed under request")
		} else {
			// 服务异常关闭
			log.Error("Server closed unexpected", "err", err)
		}
	}
}
