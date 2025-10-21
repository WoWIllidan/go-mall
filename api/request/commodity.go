package request

// CommoditySearch 商品搜索请求, --用Gin的BindQuery把URL参数绑定到结构体
type CommoditySearch struct {
	Keyword string `form:"keyword" binding:"required"`
	// 下面两个参数由Pagination组件使用
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"max=100"`
}
