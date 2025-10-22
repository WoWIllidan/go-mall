package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Order struct {
	ID          int64                 `gorm:"column:id;primary_key;AUTO_INCREMENT"`                 // 订单ID
	OrderNo     string                `gorm:"column:order_no;NOT NULL"`                             // 业务支付订单号
	PayTransId  string                `gorm:"column:pay_trans_id;NOT NULL"`                         // 支付成功后，回填的支付平台交易ID
	PayType     int                   `gorm:"column:pay_type;default:0;NOT NULL"`                   // 支付类型 0-未确定 1-微信支付 2-支付宝
	UserId      int64                 `gorm:"column:user_id;NOT NULL"`                              // 用户ID
	BillMoney   int                   `gorm:"column:bill_money;default:0;NOT NULL"`                 // 订单金额（分）
	PayMoney    int                   `gorm:"column:pay_money;default:0;NOT NULL"`                  // 支付金额（分）
	PayState    int                   `gorm:"column:pay_state;default:1;NOT NULL"`                  // 1-待支付，2-支付成功，3-支付失败
	OrderStatus int                   `gorm:"column:order_status;default:0;NOT NULL"`               // 订单状态:0.待支付 1.已支付 2.配货完成 3:已出库 4.已发货 5.配送完成待客户确认 6. 已确认收货 7. 交易成功 11.用户手动关闭 12.超时未支付关闭 13.商家确认后关闭
	PaidAt      time.Time             `gorm:"column:paid_at;default:1970-01-01 00:00:00;NOT NULL"`  // 未支付时, 默认时间为1970-01-01
	IsDel       soft_delete.DeletedAt `gorm:"softDelete:flag"`                                      // 0-未删除 1-已删除
	CreatedAt   time.Time             `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 创建时间
	UpdatedAt   time.Time             `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
}

func (Order) TableName() string {
	return "orders"
}
