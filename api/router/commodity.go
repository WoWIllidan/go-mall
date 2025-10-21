package router

import (
	"github.com/WoWBytePaladin/go-mall/api/controller"
	"github.com/gin-gonic/gin"
)

func registerCommodityRoutes(rg *gin.RouterGroup) {
	// 这个路由组中的路由都以 /commodity/ 开头
	g := rg.Group("/commodity/")
	// 按层级划分的所有商品分类
	g.GET("category-hierarchy/", controller.GetCategoryHierarchy)
	// 按ParentID 查询商品分类列表
	g.GET("category/", controller.GetCategoriesWithParentId)
}
