package model

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

// DemoOrder 只是一个示例暂时放在这里, 后面会删掉
type DemoOrder struct {
	Id           int64                 `gorm:"column:id;primary_key" json:"id"`                  //自增ID
	UserId       int64                 `gorm:"column:user_id" json:"user_id"`                    //用户ID
	BillMoney    int64                 `gorm:"column:bill_money" json:"bill_money"`              //订单金额（分）
	OrderNo      string                `gorm:"column:order_no;type:varchar(32)" json:"order_no"` //订单号
	OrderGoodsId int64                 `gorm:"column:order_goods_id" json:"order_goods_id"`      //订单对应的商品ID
	State        int8                  `gorm:"column:state;default:1" json:"state"`              //1-待支付，2-支付成功，3-支付失败
	PaidAt       time.Time             `gorm:"column:paid_at" json:"paid_at"`
	IsDel        soft_delete.DeletedAt `gorm:"softDelete:flag"`
	CreatedAt    time.Time             `gorm:"column:created_at" json:"created_at"` //创建时间
	UpdatedAt    time.Time             `gorm:"column:updated_at" json:"updated_at"` //更新时间
}

func (DemoOrder) TableName() string {
	return "demo_orders"
}
