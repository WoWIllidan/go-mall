package domainservice

import (
	"context"
	"fmt"
	"time"

	"github.com/WoWBytePaladin/go-mall/common/app"
	"github.com/WoWBytePaladin/go-mall/common/enum"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/dao"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/library"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/samber/lo"
)

type OrderDomainSvc struct {
	ctx      context.Context
	orderDao *dao.OrderDao
}

func NewOrderDomainSvc(ctx context.Context) *OrderDomainSvc {
	return &OrderDomainSvc{
		ctx:      ctx,
		orderDao: dao.NewOrderDao(ctx),
	}
}

// CreateOrder 创建订单
func (ods *OrderDomainSvc) CreateOrder(items []*do.ShoppingCartItem, userAddress *do.UserAddressInfo) (*do.Order, error) {
	billInfo, err := NewCartBillChecker(items, userAddress.UserId).GetBill()
	if err != nil {
		return nil, errcode.Wrap("CreateOrderError", err)
	}
	if billInfo.OriginalTotalPrice <= 0 {
		return nil, errcode.ErrCartItemParam
	}
	order := do.OrderNew()
	order.UserId = userAddress.UserId
	order.OrderNo = util.GenOrderNo(order.UserId)
	order.BillMoney = billInfo.OriginalTotalPrice
	order.PayMoney = billInfo.TotalPrice
	order.OrderStatus = enum.OrderStatusCreated
	if err = util.CopyProperties(&order.Items, &items); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	if err = util.CopyProperties(&order.Address, &userAddress); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	// 手动开启事务
	tx := dao.DBMaster().Begin()
	panicked := true
	defer func() { // 控制事务的提交和回滚, 保证事务的完整性
		// db.Transaction 内部其实就是这么实现的
		if err != nil || panicked { // 出现error 或者 panic 都回滚事务
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// 下面的步骤如果很多可以再使用责任链模式把步骤组织起来
	// 创建订单
	err = ods.orderDao.CreateOrder(tx, order)
	if err != nil {
		return nil, err
	}
	// 删除购物车中的购买的购物项
	cartDao := dao.NewCartDao(ods.ctx)
	cartItemIds := lo.Map(items, func(item *do.ShoppingCartItem, index int) int64 {
		return item.CartItemId
	})
	err = cartDao.DeleteMultiCartItemInTx(tx, cartItemIds)
	if err != nil {
		return nil, err
	}
	// 记录Coupon使用信息 并 锁定优惠卷
	if billInfo.Coupon.CouponId > 0 {
		// couponDao.LockCoupon(tx, coupon)
	}
	// 记录满减卷使用信息
	if billInfo.Discount.DiscountId > 0 {
		// discountDao.recordDiscount(tx, discount)
	}
	// 减少订单购买商品的库存-- 会锁行记录, 把这一步放到创建订单步骤的最后, 减少行记录加锁的时间
	commodityDao := dao.NewCommodityDao(ods.ctx)
	err = commodityDao.ReduceStuckInOrderCreate(tx, order.Items)
	if err != nil {
		return nil, err
	}

	panicked = false // 这个设置别忘了, 让事务能正常提交

	return order, err
}

// GetUserOrders 查询用户订单
func (ods *OrderDomainSvc) GetUserOrders(userId int64, pagination *app.Pagination) ([]*do.Order, error) {
	offset := pagination.Offset()
	size := pagination.GetPageSize()
	// 查询用户订单
	orderModels, totalRow, err := ods.orderDao.GetUserOrders(userId, offset, size)
	if err != nil {
		return nil, errcode.Wrap("GetUserOrdersError", err)
	}
	pagination.SetTotalRows(int(totalRow))
	orders := make([]*do.Order, 0, len(orderModels))
	if err = util.CopyProperties(&orders, &orderModels); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	// 提取所有订单ID
	orderIds := lo.Map(orders, func(order *do.Order, index int) int64 {
		return order.ID
	})
	// 查询订单的地址
	ordersAddressMap, err := ods.orderDao.GetMultiOrdersAddress(orderIds)
	if err != nil {
		return nil, errcode.Wrap("GetUserOrdersError", err)
	}
	// 查询订单明细
	ordersItemMap, err := ods.orderDao.GetMultiOrdersItems(orderIds)
	if err != nil {
		return nil, errcode.Wrap("GetUserOrdersError", err)
	}

	// 填充Order中的Address和Items
	for _, order := range orders {
		order.Address = new(do.OrderAddress) // 先初始化
		if err = util.CopyProperties(order.Address, ordersAddressMap[order.ID]); err != nil {
			return nil, errcode.ErrCoverData.WithCause(err)
		}
		orderItems := ordersItemMap[order.ID]
		if err = util.CopyProperties(&order.Items, &orderItems); err != nil {
			return nil, errcode.ErrCoverData.WithCause(err)
		}
	}

	return orders, nil
}

// GetSpecifiedUserOrder 获取 orderNo 对应的用户订单详情
func (ods *OrderDomainSvc) GetSpecifiedUserOrder(orderNo string, userId int64) (*do.Order, error) {
	orderModel, err := ods.orderDao.GetOrderByNo(orderNo)
	if err != nil {
		return nil, errcode.Wrap("GetSpecifiedUserOrderError", err)
	}
	if orderModel == nil || orderModel.UserId != userId {
		return nil, errcode.ErrOrderParams
	}
	order := do.OrderNew()
	if err = util.CopyProperties(order, orderModel); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	// 订单地址信息
	orderAddress, err := ods.orderDao.GetOrderAddress(orderModel.ID)
	if err != nil {
		return nil, errcode.Wrap("GetSpecifiedUserOrderError", err)
	}
	if err = util.CopyProperties(order.Address, orderAddress); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	// 订单购物明细
	orderItems, err := ods.orderDao.GetOrderItems(orderModel.ID)
	if err != nil {
		return nil, errcode.Wrap("GetSpecifiedUserOrderError", err)
	}
	if err = util.CopyProperties(&order.Items, &orderItems); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return order, nil
}

// CancelUserOrder 用户取消订单
func (ods *OrderDomainSvc) CancelUserOrder(orderNo string, userId int64) error {
	order, err := ods.GetSpecifiedUserOrder(orderNo, userId)
	if err != nil {
		return err
	}
	if order.OrderStatus >= enum.OrderStatusPaid {
		// 已经支付, 用户不能取消订单 -- 需要申请退款
		return errcode.ErrOrderCanNotBeChanged
	}
	// 更新订单状态为用户主动取消
	err = ods.orderDao.UpdateOrderStatus(order.ID, enum.OrderStatusUserQuit)
	if err != nil {
		return errcode.Wrap("CancelOrderError", err)

	}
	//  恢复商品的库存
	commodityDao := dao.NewCommodityDao(ods.ctx)
	err = commodityDao.RecoverOrderCommodityStuck(order.Items)
	return err
}

func (ods *OrderDomainSvc) CreteOrderWxPay(orderNo string, userId int64) (payInfo *library.WxPayInvokeInfo, err error) {
	order, err := ods.GetSpecifiedUserOrder(orderNo, userId)
	if err != nil {
		return
	}
	if order.OrderStatus != enum.OrderStatusCreated { // 订单不是初始状态,不能发起支付
		err = errcode.ErrOrderParams
		return
	}
	order.PayType = enum.PayTypeWxPay          // 支付方式--微信支付
	order.OrderStatus = enum.OrderStatusUnPaid // 订单状态--待支付
	order.PayState = enum.PayStateUnPaid
	orderModel := new(model.Order)
	if err = util.CopyProperties(orderModel, order); err != nil {
		err = errcode.Wrap("CreteOrderWxPayError", err)
		return
	}
	if err = ods.orderDao.UpdateOrder(orderModel); err != nil { // 更新Order表的支付类型和状态
		err = errcode.Wrap("CreteOrderWxPayError", err)
		return
	}

	// 用userId获取对应的Openid, 这里先Mock一个
	//openId := "QsudfrhgrDYDEEA1344EF"
	//wxPayLib := library.NewWxPayLib(ods.ctx, library.WxtPayConfig{
	//	AppId:           config.App.WechatPay.AppId,
	//	MchId:           config.App.WechatPay.MchId,
	//	PrivateSerialNo: config.App.WechatPay.PrivateSerialNo,
	//	AesKey:          config.App.WechatPay.AesKey,
	//	NotifyUrl:       config.App.WechatPay.NotifyUrl,
	//})

	// payInfo, err = wxPayLib.CreateOrderPay(order, openId)
	payInfo = &library.WxPayInvokeInfo{
		AppId:     "123456",
		TimeStamp: fmt.Sprintf("%v", time.Now().Unix()),
		NonceStr:  "e61463f8efa94090b1f366cccfbbb444",
		Package:   "prepay_id=wx21201855730335ac86f8c43d1889123400",
		SignType:  "RSA",
		PaySign:   "oR9d8PuhnIc+YZ8cBHFCwfgpaK9gd7vaRvkYD7rthRAZ/X+QBhcCYL21N7cHCTUxbQ+EAt6Uy+lwSN22f5YZvI45MLko8Pfso0jm46v5hqcVwrk6uddkGuT+Cdvu4WBqDzaDjnNa5UK3GfE1Wfl2gHxIIY5lLdUgWFts17D4WuolLLkiFZV+JSHMvH7eaLdT9N5GBovBwu5yYKUR7skR8Fu+LozcSqQixnlEZUfyE55feLOQTUYzLmR9pNtPbPsu6WVhbNHMS3Ss2+AehHvz+n64GDmXxbX++IOBvm2olHu3PsOUGRwhudhVf7UcGcunXt8cqNjKNqZLhLw4jq/xDg==",
	}

	return
}

// StartOrderWxPay 把订单设置为开始支付的状态, 支付方式为微信支付
func (ods *OrderDomainSvc) StartOrderWxPay(orderNo string, userId int64) error {
	return ods.setOrderStartPay(orderNo, userId, enum.PayTypeWxPay)
}

// setOrderStartPay 把订单设置为开始支付的状态
func (ods *OrderDomainSvc) setOrderStartPay(orderNo string, userId int64, payType int) error {
	order, err := ods.GetSpecifiedUserOrder(orderNo, userId)
	if err != nil {
		return err
	}
	if order.OrderStatus != enum.OrderStatusCreated { // 订单不是初始状态,不能发起支付
		err = errcode.ErrOrderParams
		return err
	}
	order.PayType = payType
	order.OrderStatus = enum.OrderStatusUnPaid // 订单状态--待支付
	order.PayState = enum.PayStateUnPaid
	orderModel := new(model.Order)
	if err = util.CopyProperties(orderModel, order); err != nil {
		err = errcode.Wrap("CreteOrderWxPayError", err)
		return err
	}
	if err = ods.orderDao.UpdateOrder(orderModel); err != nil { // 更新Order表的支付类型和状态
		err = errcode.Wrap("CreteOrderWxPayError", err)
		return err
	}

	return nil
}
