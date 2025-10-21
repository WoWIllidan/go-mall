package do

import (
	"time"
)

type ShoppingCartItem struct {
	CartItemId            int64  // 购物项ID
	UserId                int64  // 用户ID
	CommodityId           int64  // 商品ID
	CommodityName         string // 商品名称
	CommodityImg          string // 商品图片
	CommoditySellingPrice int    // 商品售价
	CommodityNum          int    // 商品数量
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
