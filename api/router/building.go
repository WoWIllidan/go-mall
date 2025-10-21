package router

import (
	"github.com/WoWBytePaladin/go-mall/api/controller"
	"github.com/WoWBytePaladin/go-mall/common/middleware"
	"github.com/gin-gonic/gin"
)

// 存放一些项目搭建过程中验证效果用的接口的路由

func registerBuildingRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /building 开头
	g := rg.Group("/building/")
	// 测试 Ping
	g.GET("ping", controller.TestPing)
	// 测试日志文件的读取
	g.GET("config-read", controller.TestConfigRead)
	// 测试日志门面Logger的使用
	g.GET("logger-test", controller.TestLogger)
	// 测试服务的访问日志
	g.POST("access-log-test", controller.TestAccessLog)
	// 测试服务的崩溃日志
	g.GET("panic-log-test", controller.TestPanicLog)
	// 测试项目自定义的AppError 打印错误链条和错误发生位置
	g.GET("customized-error-test", controller.TestAppError)
	// 测试统一响应--返回对象数据
	g.GET("response-obj", controller.TestResponseObj)
	// 测试统一响应--返回列表和分页
	g.GET("response-list", controller.TestResponseList)
	// 测试统一响应--返回错误
	g.GET("response-error", controller.TestResponseError)
	// 测试GORM Loggeer
	g.GET("gorm-logger-test", controller.TestGormLogger)
	// 演示代码逻辑分层, 测试 Create Demo Order
	g.POST("create-demo-order", controller.TestCreateDemoOrder)
	// 测试封装的httptool
	g.GET("httptool-get-test", controller.TestForHttpToolGet)
	g.GET("httptool-post-test", controller.TestForHttpToolPost)
	// 测试Token认证中间件
	g.GET("token-auth-test", middleware.AuthUser(), controller.TestAuthToken)
	// 初始化商品分类测试数据
	g.GET("init-category-data", controller.InitCategoryTestData)
}
