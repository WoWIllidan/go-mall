package request

type DemoOrderCreate struct {
	UserId    int64 `json:"user_id"`
	BillMoney int64 `json:"bill_money" binding:"required"`
	// 这个字段演示的时候因为没创建订单快照表所以不写库
	OrderGoodsId int64 `json:"order_goods_id" binding:"required"`
}
