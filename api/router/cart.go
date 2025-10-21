package router

import (
	"github.com/WoWBytePaladin/go-mall/api/controller"
	"github.com/WoWBytePaladin/go-mall/common/middleware"
	"github.com/gin-gonic/gin"
)

// 存放购物车模块的路由

func registerCartRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /cart/ 开头, 并且都需要身份验证
	g := rg.Group("/cart/")
	g.Use(middleware.AuthUser())
	// 添加到购物车
	g.POST("add-item", controller.AddCartItem)
	// 修改购物车中的商品数量
	g.PATCH("update-item", controller.UpdateCartItem)
	// 用户购物车中的购物项列表
	g.GET("item/", controller.UserCartItems)
	// 删除购物项
	g.DELETE("/item/:item_id", controller.DeleteUserCartItem)
	// 查看购物项账单 -- 确认下单前用来显示商品和支付金额明细
	g.GET("/item/check-bill", controller.CheckCartItemBill)
}
