package model

import (
	"time"
)

// 订单收货信息 -- 订单快照

type OrderAddress struct {
	ID            int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 订单关联收货信息的id
	OrderId       int64     `gorm:"column:order_id;NOT NULL"`                             // 订单id
	UserName      string    `gorm:"column:user_name;NOT NULL"`                            // 收货人姓名
	UserPhone     string    `gorm:"column:user_phone;NOT NULL"`                           // 收货人手机号
	ProvinceName  string    `gorm:"column:province_name;NOT NULL"`                        // 省
	CityName      string    `gorm:"column:city_name;NOT NULL"`                            // 城
	RegionName    string    `gorm:"column:region_name;NOT NULL"`                          // 区
	DetailAddress string    `gorm:"column:detail_address;NOT NULL"`                       // 收件详细地址(街道/楼宇/单元)
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

func (m *OrderAddress) TableName() string {
	return "order_address"
}
