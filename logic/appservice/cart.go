package appservice

import (
	"context"
	"github.com/WoWBytePaladin/go-mall/api/reply"
	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/WoWBytePaladin/go-mall/logic/domainservice"
	"github.com/samber/lo"
)

type CartAppSvc struct {
	ctx           context.Context
	cartDomainSvc *domainservice.CartDomainSvc
}

func NewCartAppSvc(ctx context.Context) *CartAppSvc {
	return &CartAppSvc{
		ctx:           ctx,
		cartDomainSvc: domainservice.NewCartDomainSvc(ctx),
	}
}

// AddCartItem 添加商品到购物车
func (cas *CartAppSvc) AddCartItem(request *request.AddCartItem, userId int64) error {
	commodityDomainSvc := domainservice.NewCommodityDomainSvc(cas.ctx)
	commodityInfo := commodityDomainSvc.GetCommodityInfo(request.CommodityId)
	if commodityInfo == nil || commodityInfo.ID == 0 { // 商品不存在
		return errcode.ErrCommodityNotExists
	}
	if commodityInfo.StockNum < request.CommodityNum {
		// 先初步判断库存是否充足, 下单时需要重新用当前读判断库存
		return errcode.ErrCommodityStockOut
	}

	shoppCartItem := new(do.ShoppingCartItem)
	err := util.CopyProperties(shoppCartItem, request)
	if err != nil {
		logger.New(cas.ctx).Error(errcode.ErrCoverData.Msg(), "err", err)
		return err
	}
	shoppCartItem.UserId = userId

	return cas.cartDomainSvc.CartAddItem(shoppCartItem)
}

// UpdateCartItem 更新购物项
func (cas *CartAppSvc) UpdateCartItem(request *request.CartItemUpdate, userId int64) error {
	return cas.cartDomainSvc.CartUpdateItem(request, userId)
}

// GetUserCartItems 获取用户购物车中的购物项
func (cas *CartAppSvc) GetUserCartItems(userId int64) ([]*reply.CartItem, error) {
	cartItems, err := cas.cartDomainSvc.GetUserCartItems(userId)
	if err != nil {
		return nil, err
	}
	replyCartItems := make([]*reply.CartItem, 0, len(cartItems))
	err = util.CopyProperties(&replyCartItems, &cartItems)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return replyCartItems, nil
}

// DeleteUserCartItem 删除用户购物车中的购物项, 只能删除单个购物项
func (cas *CartAppSvc) DeleteUserCartItem(cartItemId, userId int64) error {
	return cas.cartDomainSvc.DeleteUserCartItem(cartItemId, userId)
}

// CheckCartItemBill 查看购物项账单
func (cas *CartAppSvc) CheckCartItemBill(cartItemIds []int64, userId int64) (*reply.CheckedCartItemBill, error) {
	checkedCartItems, err := cas.cartDomainSvc.GetCheckedCartItems(cartItemIds, userId)
	if err != nil {
		return nil, err
	}

	// 计算总价
	totalPrice := lo.Reduce(checkedCartItems, func(agg int, item *do.ShoppingCartItem, index int) int {
		return agg + item.CommoditySellingPrice*item.CommodityNum
	}, 0)

	replyBill := new(reply.CheckedCartItemBill)
	err = util.CopyProperties(&replyBill.Items, checkedCartItems)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	replyBill.TotalPrice = totalPrice
	return replyBill, nil
}

// CheckCartItemBillV2 V2版购物项账单, 支持满减、优惠卷、会员价
func (cas *CartAppSvc) CheckCartItemBillV2(cartItemIds []int64, userId int64) (*reply.CheckedCartItemBillV2, error) {
	checkedCartItems, err := cas.cartDomainSvc.GetCheckedCartItems(cartItemIds, userId)
	if err != nil {
		return nil, err
	}
	billChecker := domainservice.NewCartBillChecker(checkedCartItems, userId)
	billInfo, err := billChecker.GetBill()
	if err != nil {
		return nil, err
	}
	replyBillInfo := new(reply.CheckedCartItemBillV2)
	if err = util.CopyProperties(&replyBillInfo.Items, checkedCartItems); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	if err = util.CopyProperties(&replyBillInfo.BillDetail, &billInfo); err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	return replyBillInfo, nil
}
