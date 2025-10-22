package model

import (
	"time"
)

// 订单购物项 -- 订单快照

type OrderItem struct {
	ID                    int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 订单关联购物项主键id
	OrderId               int64     `gorm:"column:order_id;NOT NULL"`                             // 订单主键id
	CommodityId           int64     `gorm:"column:commodity_id;NOT NULL"`                         // 关联的商品id
	CommodityName         string    `gorm:"column:commodity_name;NOT NULL"`                       // 下单时商品的名称(订单快照)
	CommodityImg          string    `gorm:"column:commodity_img;NOT NULL"`                        // 下单时商品的主图(订单快照)
	CommoditySellingPrice int       `gorm:"column:commodity_selling_price;default:0;NOT NULL"`    // 下单时商品的价格(订单快照)
	CommodityNum          int       `gorm:"column:commodity_num;default:1;NOT NULL"`              // 数量(订单快照)
	CreatedAt             time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt             time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

func (OrderItem) TableName() string {
	return "order_items"
}
