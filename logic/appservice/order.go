package appservice

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/api/reply"
	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/app"
	"github.com/WoWBytePaladin/go-mall/common/enum"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/logic/domainservice"
)

type OrderAppSvc struct {
	ctx            context.Context
	orderDomainSvc *domainservice.OrderDomainSvc
}

func NewOrderAppSvc(ctx context.Context) *OrderAppSvc {
	return &OrderAppSvc{
		ctx:            ctx,
		orderDomainSvc: domainservice.NewOrderDomainSvc(ctx),
	}
}

// CreateOrder 创建订单
func (oas *OrderAppSvc) CreateOrder(orderRequest *request.OrderCreate, userId int64) (*reply.OrderCreateReply, error) {
	cartDomainSvc := domainservice.NewCartDomainSvc(oas.ctx)
	cartItems, err := cartDomainSvc.GetCheckedCartItems(orderRequest.CartItemIdList, userId)
	if err != nil {
		return nil, err
	}
	userDomainSvc := domainservice.NewUserDomainSvc(oas.ctx)
	address, err := userDomainSvc.GetUserSingleAddress(userId, orderRequest.UserAddressId)
	if err != nil {
		return nil, err
	}

	order, err := oas.orderDomainSvc.CreateOrder(cartItems, address)
	if err != nil {
		return nil, err
	}
	orderReply := new(reply.OrderCreateReply)
	orderReply.OrderNo = order.OrderNo
	return orderReply, nil
}

// GetUserOrders 查询用户订单
func (oas *OrderAppSvc) GetUserOrders(userId int64, pagination *app.Pagination) ([]*reply.Order, error) {
	orders, err := oas.orderDomainSvc.GetUserOrders(userId, pagination)
	if err != nil {
		return nil, err
	}
	replyOrders := make([]*reply.Order, 0, len(orders))

	if err = util.CopyProperties(&replyOrders, &orders); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	for _, replyOrder := range replyOrders {
		// 订单的前台状态
		replyOrder.FrontStatus = enum.OrderFrontStatus[replyOrder.OrderStatus]
		replyOrder.Address.UserName = util.MaskRealName(replyOrder.Address.UserName)
		replyOrder.Address.UserPhone = util.MaskPhone(replyOrder.Address.UserPhone)
	}
	return replyOrders, nil
}

// GetOrderInfo 订单详情
func (oas *OrderAppSvc) GetOrderInfo(orderNo string, userId int64) (*reply.Order, error) {
	order, err := oas.orderDomainSvc.GetSpecifiedUserOrder(orderNo, userId)
	if err != nil {
		return nil, err
	}

	replyOrder := new(reply.Order)
	if err = util.CopyProperties(replyOrder, order); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	// 订单的前台状态
	replyOrder.FrontStatus = enum.OrderFrontStatus[replyOrder.OrderStatus]
	// 敏感信息脱敏
	replyOrder.Address.UserName = util.MaskRealName(replyOrder.Address.UserName)
	replyOrder.Address.UserPhone = util.MaskPhone(replyOrder.Address.UserPhone)

	return replyOrder, nil
}

// CancelOrder 取消订单
func (oas *OrderAppSvc) CancelOrder(orderNo string, userId int64) error {
	return oas.orderDomainSvc.CancelUserOrder(orderNo, userId)
}

// OrderCreatePay 订单发起支付
func (oas *OrderAppSvc) OrderCreatePay(payRequest *request.OrderPayCreate, userId int64) (replyData interface{}, err error) {
	switch payRequest.PayType {
	case enum.PayTypeWxPay: // 使用微信支付
		payInfo, err := oas.orderDomainSvc.CreteOrderWxPay(payRequest.OrderNo, userId)
		return payInfo, err
	default:
		err = errcode.ErrParams
	}

	return
}
