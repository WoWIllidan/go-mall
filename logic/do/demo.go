package do

import (
	"time"
)

// 演示DEMO, 后期使用时删掉

type DemoOrder struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"userId"`
	BillMoney    int64     `json:"billMoney"`
	OrderNo      string    `json:"orderNo"`
	OrderGoodsId int64     `json:"orderGoodsId"`
	State        int8      `json:"state"`
	PaidAt       time.Time `json:"paidAt"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
