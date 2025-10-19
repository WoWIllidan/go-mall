package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID        int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 用户ID
	Nickname  string                `gorm:"column:nickname;NOT NULL"`                             // 用户昵称
	LoginName string                `gorm:"column:login_name;NOT NULL"`                           // 登录时使用的用户名
	Password  string                `gorm:"column:password;NOT NULL"`                             // bcrypt加密的登录密码
	Verified  int                   `gorm:"column:verified;default:0;NOT NULL"`                   // 验证状态 0-未验证 1-已验证
	Avatar    string                `gorm:"column:avatar;NOT NULL"`                               // 用户头像
	Slogan    string                `gorm:"column:slogan;NOT NULL"`                               // 个性签名
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 删除状态 0-未删除 1-已删除
	IsBlocked int                   `gorm:"column:is_blocked;default:0;NOT NULL"`                 // 禁用状态 0-正常 1-已禁用
	CreatedAt time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

func (u *User) TableName() string {
	return "users"
}
