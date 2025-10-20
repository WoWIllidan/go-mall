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
	// 新增用户收货地址信息
	g.POST("address", middleware.AuthUser(), controller.AddUserAddress)
	// 查询用户所有的收货地址信息
	g.GET("address/", middleware.AuthUser(), controller.GetUserAddresses)
	// 查询用户单个收货地址信息
	g.GET("address/:address_id", middleware.AuthUser(), controller.GetUserAddress)
	// 修改用户的地址
	g.PATCH("address/:address_id", middleware.AuthUser(), controller.UpdateUserAddress)
	// 删除用户的单条地址信息
	g.DELETE("address/:address_id", middleware.AuthUser(), controller.DeleteUserAddress)
}
