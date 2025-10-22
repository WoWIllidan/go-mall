package do

import (
	"time"
)

type Order struct {
	ID          int64
	OrderNo     string
	PayTransId  string
	PayType     int
	UserId      int64
	BillMoney   int
	PayMoney    int
	PayState    int
	OrderStatus int
	Address     *OrderAddress
	Items       []*OrderItem
	PaidAt      time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderAddress struct {
	//ID            int64  //领域对象里里 OrderAddress不需要ID，它依附再Order对象上，
	// 同时有ID，Copy的时候会把UserAddress的ID复制到OrderAddress的ID上，写orderAddress表时会出现主键冲突
	OrderId       int64
	UserName      string
	UserPhone     string
	ProvinceName  string
	CityName      string
	RegionName    string
	DetailAddress string
}

type OrderItem struct {
	OrderId               int64
	CommodityId           int64
	CommodityName         string
	CommodityImg          string
	CommoditySellingPrice int
	CommodityNum          int
}

func OrderNew() *Order {
	order := new(Order)
	order.Address = new(OrderAddress) // 内嵌的Pointer字段不自己初始化会是 nil, 无法用 util.CopyProperties 来拷贝属性值
	return order
}
