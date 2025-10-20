package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// UserAddress 用户收货信息表
type UserAddress struct {
	ID            int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 收货信息ID
	UserId        int64                 `gorm:"column:user_id;NOT NULL"`                              // 用户ID
	UserName      string                `gorm:"column:user_name;NOT NULL"`                            // 收货人姓名
	UserPhone     string                `gorm:"column:user_phone;NOT NULL"`                           // 收货人手机号
	Default       int                   `gorm:"column:default;default:0;NOT NULL"`                    // 是否为默认收货信息 0-非默认 1-是默认
	ProvinceName  string                `gorm:"column:province_name;NOT NULL"`                        // 省
	CityName      string                `gorm:"column:city_name;NOT NULL"`                            // 城
	RegionName    string                `gorm:"column:region_name;NOT NULL"`                          // 区/县
	DetailAddress string                `gorm:"column:detail_address;NOT NULL"`                       // 收件详细地址(街道/楼宇/单元)
	IsDel         soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 删除状态 0-未删除 1-已删除
	CreatedAt     time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 添加时间
	UpdatedAt     time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 修改时间
}

func (m *UserAddress) TableName() string {
	return "user_address"
}
