package domainservice

import (
	"context"
	"errors"

	"github.com/WoWBytePaladin/go-mall/common/enum"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/config"
	"github.com/WoWBytePaladin/go-mall/library"
	"github.com/WoWBytePaladin/go-mall/logic/do"
)

// OrderPayTemplateContract 订单支付的模版--对订单支付执行过程的抽象, 模版方法中决定流程步骤的执行顺序
type OrderPayTemplateContract interface {
	CreateOrderPay() (interface{}, error) // 模版方法
	OrderPayHandlerContract
}

// OrderPayHandlerContract 订单支付的处理器接口--对订单支付各个主要步骤的抽象
type OrderPayHandlerContract interface {
	// CheckRepetition 防重校验
	CheckRepetition() error
	// ValidateOrder 检验Order参数是否符合预期
	ValidateOrder() error
	// LoadPayAndUserConfig 加载支付配置和支付平台需要的一些用户信息--比如微信的openID
	LoadPayAndUserConfig() error
	// LoadOrderPayStrategy 加载订单支付策略
	LoadOrderPayStrategy() error
	// HandleOrderPay 发起支付
	HandleOrderPay() (interface{}, error)
}

// OrderPayStrategyContract 支付策略接口--接口的实现类去真正发起对支付平台的调用
type OrderPayStrategyContract interface {
	// CreatePay 实现支付策略中有关创建支付的逻辑
	// 可以根据策略本身决定去怎么实现 创建订单|设置签约码|......
	CreatePay(ctx context.Context, order *do.Order, payConfig *OrderPayConfig) (interface{}, error)
}

// OrderPayTemplate 是实现 OrderPayTemplateContract 的抽象类型
// 只实现模版方法 CreateOrderPay , 在其中规定发起支付要执行的步骤的顺序
// 具体的执行步骤的实现放到 OrderPayHandlerContract 的实现中
type OrderPayTemplate struct {
	OrderPayHandlerContract
}

func (template OrderPayTemplate) CreateOrderPay() (interface{}, error) {
	// 防止用户端重复操作
	if err := template.CheckRepetition(); err != nil {
		return nil, err
	}
	// 校验参数是否符合预期
	if err := template.ValidateOrder(); err != nil {
		return nil, err
	}

	// 加载支付配置和支付平台需要的一些用户信息
	if err := template.LoadPayAndUserConfig(); err != nil {
		return nil, err
	}
	// 加载支付策略--例如微信的小程序支付
	if err := template.LoadOrderPayStrategy(); err != nil {
		return nil, err
	}
	response, err := template.HandleOrderPay()
	if err != nil {
		return nil, err
	}

	return response, nil
}

type OrderPayConfig struct {
	PayUserId   int64
	WxOpenId    string
	WxPayConfig *library.WxtPayConfig
	//AliPayConfig
}

// CommonOrderPayHandler 支付处理的通用类, 只实现参数校验这样的每个支付方式都需要做的通用操作。
// 其他操作都交由具体的支付方式类去自己实现覆盖这里的默认实现
type CommonOrderPayHandler struct {
	ctx       context.Context
	Scene     string // 支付场景 H5 、app、小程序 jsapi(公众号、线下、PC网页) 等 -- 对应支付平台不同支付场景
	UserId    int64
	OrderNo   string // 业务订单号
	Order     *do.Order
	PayConfig *OrderPayConfig
	// 这里还可以继续定义AliPayConfig 字段

	PayStrategy OrderPayStrategyContract // 支付策略
}

func (handler *CommonOrderPayHandler) CheckRepetition() error {
	// 根据自己的业务量, 用 Redis 或者更高级的方式做防重校验
	return nil
}

func (handler *CommonOrderPayHandler) ValidateOrder() error {
	order, err := NewOrderDomainSvc(handler.ctx).GetSpecifiedUserOrder(handler.OrderNo, handler.UserId)
	if err != nil {
		return err
	}
	if order.OrderStatus > enum.OrderStatusCreated {
		return errcode.ErrOrderParams // 订单状态不对, 不能发起支付
	}

	handler.Order = order
	return nil
}

func (handler *CommonOrderPayHandler) LoadPayAndUserConfig() error {
	// 交给具体的支付方式类型去实现
	return nil
}

func (handler *CommonOrderPayHandler) LoadOrderPayStrategy() error {
	//handler.PayStrategy = new(...)
	// 留给具体的支付方式类型去实现
	return nil
}

func (handler *CommonOrderPayHandler) HandleOrderPay() (interface{}, error) {
	return handler.PayStrategy.CreatePay(handler.ctx, handler.Order, handler.PayConfig)
}

// WxOrderPayHandler 微信订单支付处理类
type WxOrderPayHandler struct {
	CommonOrderPayHandler
}

func (wxHandler *WxOrderPayHandler) LoadPayAndUserConfig() error {
	wxHandler.PayConfig.WxPayConfig = &library.WxtPayConfig{
		AppId:           config.App.WechatPay.AppId,
		MchId:           config.App.WechatPay.MchId,
		PrivateSerialNo: config.App.WechatPay.PrivateSerialNo,
		AesKey:          config.App.WechatPay.AesKey,
		NotifyUrl:       config.App.WechatPay.NotifyUrl,
	}
	wxHandler.PayConfig.PayUserId = wxHandler.UserId
	// 用userId获取对应的Openid, 这里先Mock一个
	// xxx.GetUserOpenId(wxHandler.userId)
	openId := "QsudfrhgrDYDEEA1344EF"
	wxHandler.PayConfig.WxOpenId = openId
	return nil
}

func (wxHandler *WxOrderPayHandler) LoadOrderPayStrategy() error {
	switch wxHandler.Scene {
	case "app": // app 支付
		return errcode.ErrOrderUnsupportedPayScene
	// 加载app的支付策略类
	case "jsapi": // 网页支付
		// 加载封装了微信支付JSAPI的策略类
		wxHandler.PayStrategy = new(WxJSPayStrategy)
	default:
		return errcode.ErrOrderParams.WithCause(errors.New("unsupported platform"))
	}

	return nil
}

// WxJSPayStrategy 微信JSAPI 支付接口实现
type WxJSPayStrategy struct {
}

func (strategy *WxJSPayStrategy) CreatePay(ctx context.Context, order *do.Order, payConfig *OrderPayConfig) (interface{}, error) {
	ods := NewOrderDomainSvc(ctx)
	if err := ods.StartOrderWxPay(order.OrderNo, order.UserId); err != nil {
		return nil, err
	}
	wpl := library.NewWxPayLib(ctx, *payConfig.WxPayConfig)
	reply, err := wpl.CreateOrderPay(order, payConfig.WxOpenId)
	if err != nil {
		err = errcode.Wrap("WxJSPayStrategyCreatePayError", err)
	}
	return reply, err
}

// type WxAppPayStrategy
// ......

// NewOrderPayTemplate
// 创建订单支付模版的工厂方法
// @param ctx
// @param userId
// @param orderNo
// @param payScene 支付场景 app h5 jsapi min-app...
// @param payType 支付类型  微信支付｜支付宝 ｜ ...
func NewOrderPayTemplate(ctx context.Context, userId int64, orderNo, payScene string, payType int) *OrderPayTemplate {
	payTemplate := new(OrderPayTemplate)
	switch payType {
	case enum.PayTypeWxPay: // 微信支付
		payHandler := new(WxOrderPayHandler)
		payHandler.ctx = ctx
		payHandler.UserId = userId
		payHandler.OrderNo = orderNo
		payHandler.Scene = payScene
		payTemplate.OrderPayHandlerContract = payHandler

	}

	return payTemplate
}
