package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Commodity struct {
	ID            int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 商品表主键id
	Name          string                `gorm:"column:name;NOT NULL"`                                 // 商品名
	Intro         string                `gorm:"column:intro;NOT NULL"`                                // 商品简介
	CategoryId    int64                 `gorm:"column:category_id;default:0;NOT NULL"`                // 关联分类id
	CoverImg      string                `gorm:"column:cover_img;NOT NULL"`                            // 商品封面图
	Images        string                `gorm:"column:images;NOT NULL"`                               // 商品细节图
	DetailContent string                `gorm:"column:detail_content;NOT NULL"`                       // 商品详情
	OriginalPrice int                   `gorm:"column:original_price;default:1;NOT NULL"`             // 商品原价
	SellingPrice  int                   `gorm:"column:selling_price;default:1;NOT NULL"`              // 商品售价
	StockNum      int                   `gorm:"column:stock_num;default:0;NOT NULL"`                  // 商品库存数量
	Tag           string                `gorm:"column:tag;NOT NULL"`                                  // 商品标签
	SellStatus    int                   `gorm:"column:sell_status;default:1;NOT NULL"`                // 商品上架状态 1-上架  2-下架
	IsDel         soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 删除标识字段(0-未删除 1-已删除)
	CreatedAt     time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt     time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

func (m *Commodity) TableName() string {
	return "commodities"
}
