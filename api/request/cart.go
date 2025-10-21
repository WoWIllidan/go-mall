package request

// AddCartItem 添加到购物车
type AddCartItem struct {
	CommodityId  int64 `json:"commodity_id" binding:"required"`
	CommodityNum int   `json:"commodity_num" binding:"required,min=1,max=5"` // 一个商品往购物车里一次性最多放5个
}

type CartItemUpdate struct {
	ItemId       int64 `json:"item_id" binding:"required"`
	CommodityNum int   `json:"commodity_num" binding:"required,min=1,max=6"`
}
