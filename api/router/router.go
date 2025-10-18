package router

import (
	"github.com/WoWBytePaladin/go-mall/common/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	// use global middlewares
	engine.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	routeGroup := engine.Group("")
	registerBuildingRoutes(routeGroup)
}
