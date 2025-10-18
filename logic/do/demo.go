package do

import (
	"time"
)

// 演示DEMO, 后期使用时删掉

type DemoOrder struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	BillMoney int64     `json:"bill_money"`
	OrderNo   string    `json:"order_no"`
	State     int8      `json:"state"`
	IsDel     uint      `json:"is_del"`
	PaidAt    time.Time `json:"paid_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
