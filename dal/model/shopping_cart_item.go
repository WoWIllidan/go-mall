package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// ShoppingCartItem 购物车明细表
type ShoppingCartItem struct {
	CartItemId   int64                 `gorm:"column:cart_item_id;primary_key;AUTO_INCREMENT"`       // 购物项主键id
	UserId       int64                 `gorm:"column:user_id;NOT NULL"`                              // 用户主键id
	CommodityId  int64                 `gorm:"column:commodity_id;NOT NULL"`                         // 关联商品id
	CommodityNum int                   `gorm:"column:commodity_num;default:1;NOT NULL"`              // 商品数量
	IsDel        soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 删除(0-未删除 1-已删除)
	CreatedAt    time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt    time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

func (ShoppingCartItem) TableName() string {
	return "shopping_cart_items"
}
