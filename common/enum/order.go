package enum

const (
	PayStateNotInitiated = iota
	PayStateUnPaid
	PayStatePaid
	PayStatePayFailed
)

const (
	OrderStatusCreated        = iota // 已创建
	OrderStatusUnPaid                // 待支付
	OrderStatusPaid                  // 已支付
	OrderStatusChecked               // 检货完成
	OrderStatusShipped               // 已发货
	OrderStatusOnDelivery            // 配送中 -- 快递员上门送货中
	OrderStatusDelivered             // 已送达
	OrderStatusConfirmReceipt        // 已确认收货
	OrderStatusCompleted             // 订单完成
	OrderStatusUserQuit              // 用户取消
	OrderStatusUnpaidClose           // 超时未支付
	OrderStatusMerchantClose         // 商家关闭订单
)

// OrderFrontStatus 用户在前台看到的订单状态
var OrderFrontStatus = map[int]string{
	OrderStatusCreated:        "待付款",
	OrderStatusUnPaid:         "待付款",
	OrderStatusPaid:           "待发货",
	OrderStatusChecked:        "待发货",
	OrderStatusShipped:        "待收货",
	OrderStatusOnDelivery:     "待收货",
	OrderStatusDelivered:      "待收货",
	OrderStatusConfirmReceipt: "待评价",
	OrderStatusCompleted:      "已完成",
	OrderStatusUserQuit:       "已取消",
	OrderStatusUnpaidClose:    "已取消",
	OrderStatusMerchantClose:  "已取消",
}
