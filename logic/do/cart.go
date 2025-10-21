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

type CartBillInfo struct {
	Coupon struct { // 可用的优惠券
		CouponId      int64
		CouponName    string
		DiscountMoney int
		Threshold     int // 使用门槛, 比如满1000 可用

	}
	Discount struct { // 可用的满减活动券
		DiscountId    int64
		DiscountName  string
		DiscountMoney int
		Threshold     int // 使用门槛, 比如满1000 可用
	}
	VipDiscountMoney   int // VIP减免的金额
	OriginalTotalPrice int // 减免、优惠前的总金额
	TotalPrice         int // 实际要支付的总金额
}
