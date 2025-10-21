package reply

type CartItem struct {
	CartItemId            int64  `json:"cart_item_id"`
	UserId                int64  `json:"user_id"`
	CommodityId           int64  `json:"commodity_id"`
	CommodityNum          int    `json:"commodity_num"`
	CommodityName         string `json:"commodity_name"`                 // 商品名称
	CommodityImg          string `json:"commodity_img"`                  // 商品图片
	CommoditySellingPrice int    `json:"commodity_selling_price"`        // 商品售价
	AddCartAt             string `json:"add_cart_at" copier:"CreatedAt"` //购物车添加时间,  把Do的CreatedAt字段用copier映射到这里
}

type CheckedCartItemBill struct {
	Items      []*CartItem `json:"items"`
	TotalPrice int         `json:"total_price"` // 总价
}
