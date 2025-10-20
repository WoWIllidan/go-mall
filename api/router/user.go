package router

import (
	"github.com/WoWBytePaladin/go-mall/api/controller"
	"github.com/WoWBytePaladin/go-mall/common/middleware"
	"github.com/gin-gonic/gin"
)

// 存放User模块的路由

func registerUserRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /user 开头
	g := rg.Group("/user/")
	// 刷新Token
	g.GET("token/refresh", controller.RefreshUserToken)
	// 注册用户
	g.POST("register", controller.RegisterUser)
	// 登录
	g.POST("login", controller.LoginUser)
	// 登出用户
	g.DELETE("logout", middleware.AuthUser(), controller.LogoutUser)
	// 申请重置密码
	g.POST("password/apply-reset", controller.PasswordResetApply)
	// 重置密码
	g.POST("password/reset", controller.PasswordReset)
	// 用户基本信息
	g.GET("info", middleware.AuthUser(), controller.UserInfo)
	// 更新用户基本信息
	g.PATCH("info", middleware.AuthUser(), controller.UpdateUserInfo)
}
