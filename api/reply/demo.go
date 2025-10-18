package reply

type DemoOrder struct {
	UserId    int64  `json:"user_id"`
	BillMoney int64  `json:"bill_money"`
	OrderNo   string `json:"order_no"`
	State     int8   `json:"state"`
	PaidAt    string `json:"paid_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
