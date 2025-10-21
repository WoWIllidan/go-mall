package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// CommodityCategory 商品分类表
type CommodityCategory struct {
	ID        int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 分类id
	Level     int                   `gorm:"column:level;default:0;NOT NULL"`                      // 分类级别(1-一级分类 2-二级分类 3-三级分类)
	ParentId  int64                 `gorm:"column:parent_id;default:0;NOT NULL"`                  // 父分类id
	Name      string                `gorm:"column:name;NOT NULL"`                                 // 分类名称
	IconImg   string                `gorm:"column:icon_img;NOT NULL"`                             // 分类的图标
	Rank      int                   `gorm:"column:rank;default:0;NOT NULL"`                       // 排序值(字段越大越靠前)
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 删除标识字段(0-未删除 1-已删除)
	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}
