package reply

type OrderCreateReply struct {
	OrderNo string `json:"order_no"`
}

type Order struct {
	OrderNo     string `json:"order_no"`
	PayTransId  string `json:"pay_trans_id"`
	PayType     int    `json:"pay_type"`
	BillMoney   int    `json:"bill_money"`
	PayMoney    int    `json:"pay_money"`
	PayState    int    `json:"pay_state"`
	OrderStatus int    `json:"-"`
	FrontStatus string `json:"status"`
	Address     struct {
		UserName      string `json:"user_name"`
		UserPhone     string `json:"user_phone"`
		ProvinceName  string `json:"province_name"`
		CityName      string `json:"city_name"`
		RegionName    string `json:"region_name"`
		DetailAddress string `json:"detail_address"`
	} `json:"address,omitempty"`
	Items []struct {
		CommodityId           int64  `json:"commodity_id"`
		CommodityName         string `json:"commodity_name"`
		CommodityImg          string `json:"commodity_img"`
		CommoditySellingPrice int    `json:"commodity_selling_price"`
		CommodityNum          int    `json:"commodity_num"`
	} `json:"items,omitempty"`
	CreatedAt string `json:"created_at"`
}
