package controller

import (
	"errors"
	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/app"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/logic/appservice"
	"github.com/gin-gonic/gin"
)

// OrderCreate 创建订单
func OrderCreate(c *gin.Context) {
	request := new(request.OrderCreate)
	if err := c.ShouldBindJSON(request); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	orderAppSvc := appservice.NewOrderAppSvc(c)
	reply, err := orderAppSvc.CreateOrder(request, c.GetInt64("userId"))
	if err != nil {
		if errors.Is(err, errcode.ErrCartItemParam) {
			app.NewResponse(c).Error(errcode.ErrCartItemParam)
		} else if errors.Is(err, errcode.ErrCartWrongUser) {
			app.NewResponse(c).Error(errcode.ErrCartWrongUser)
		} else if errors.Is(err, errcode.ErrCommodityStockOut) {
			app.NewResponse(c).Error(errcode.ErrCommodityStockOut.WithCause(err))
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(reply)
}

// UserOrders 用户订单列表
func UserOrders(c *gin.Context) {
	pagination := app.NewPagination(c)
	orderAppSvc := appservice.NewOrderAppSvc(c)
	replyOrders, err := orderAppSvc.GetUserOrders(c.GetInt64("userId"), pagination)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
	}

	app.NewResponse(c).SetPagination(pagination).Success(replyOrders)
}

// OrderInfo 订单详情
func OrderInfo(c *gin.Context) {
	orderNo := c.Param("order_no")
	orderAppSvc := appservice.NewOrderAppSvc(c)
	replyOrder, err := orderAppSvc.GetOrderInfo(orderNo, c.GetInt64("userId"))
	if err != nil {
		if errors.Is(err, errcode.ErrOrderParams) {
			app.NewResponse(c).Error(errcode.ErrOrderParams)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(replyOrder)
}

// OrderCancel 用户主动取消订单
func OrderCancel(c *gin.Context) {
	orderNo := c.Param("order_no")
	orderAppSvc := appservice.NewOrderAppSvc(c)
	err := orderAppSvc.CancelOrder(orderNo, c.GetInt64("userId"))
	if err != nil {
		if errors.Is(err, errcode.ErrOrderParams) {
			app.NewResponse(c).Error(errcode.ErrOrderParams)
		} else if errors.Is(err, errcode.ErrOrderCanNotBeChanged) {
			app.NewResponse(c).Error(errcode.ErrOrderCanNotBeChanged)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).SuccessOk()
}

// CreateOrderPay 订单发起支付
func CreateOrderPay(c *gin.Context) {
	request := new(request.OrderPayCreate)
	if err := c.ShouldBindJSON(request); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	orderAppSvc := appservice.NewOrderAppSvc(c)
	reply, err := orderAppSvc.OrderCreatePay(request, c.GetInt64("userId"))
	if err != nil {
		if errors.Is(err, errcode.ErrOrderParams) {
			app.NewResponse(c).Error(errcode.ErrOrderParams)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(reply)
}
