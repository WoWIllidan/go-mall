package domainservice

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/dao"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/samber/lo"
)

type CartDomainSvc struct {
	ctx     context.Context
	cartDao *dao.CartDao
}

func NewCartDomainSvc(ctx context.Context) *CartDomainSvc {
	return &CartDomainSvc{
		ctx:     ctx,
		cartDao: dao.NewCartDao(ctx),
	}
}

// CartAddItem 购物车添加商品
func (cds *CartDomainSvc) CartAddItem(cartItem *do.ShoppingCartItem) error {
	cartItemModel, err := cds.cartDao.GetUserCartItemWithCommodityId(cartItem.UserId, cartItem.CommodityId)
	if err != nil {
		return errcode.Wrap("CartAddItemError", err)
	}
	if cartItemModel != nil && cartItemModel.CartItemId != 0 {
		// 添加购物车中已经存在的项目, 商品数量累加更新(可根据产品逻辑限制一个商品的最大数)
		cartItemModel.CommodityNum += cartItem.CommodityNum
		return cds.cartDao.UpdateCartItem(cartItemModel)
	}

	err = cds.cartDao.AddCartItem(cartItem)
	if err != nil {
		err = errcode.Wrap("CartAddItemError", err)
	}
	return err
}

// CartUpdateItem 更改购物项
func (cds *CartDomainSvc) CartUpdateItem(request *request.CartItemUpdate, userId int64) error {
	cartItemModel, err := cds.cartDao.GetCartItemById(request.ItemId)
	if err != nil {
		err = errcode.Wrap("CartUpdateItemError", err)
		return err
	}
	if cartItemModel == nil || cartItemModel.UserId != userId {
		logger.New(cds.ctx).Error("DataMatchError", "cartItem", cartItemModel, "request", request, "requestUserId", userId)
		return errcode.ErrParams
	}
	cartItemModel.CommodityNum = request.CommodityNum
	err = cds.cartDao.UpdateCartItem(cartItemModel)
	if err != nil {
		err = errcode.Wrap("CartUpdateItemError", err)
	}

	return err
}

// GetUserCartItems 获取用户购物车里的购物项
func (cds *CartDomainSvc) GetUserCartItems(userId int64) ([]*do.ShoppingCartItem, error) {
	cartItemModels, err := cds.cartDao.GetUserCartItems(userId)
	if err != nil {
		err = errcode.Wrap("GetUserCartItemsError", err)
		return nil, err
	}
	userCartItems := make([]*do.ShoppingCartItem, 0, len(cartItemModels))
	err = util.CopyProperties(&userCartItems, &cartItemModels)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	err = cds.fillInCommodityInfo(userCartItems)
	if err != nil {
		return nil, err
	}

	return userCartItems, nil
}

// fillInCommodityInfo 为购物项填充商品信息
func (cds *CartDomainSvc) fillInCommodityInfo(cartItems []*do.ShoppingCartItem) error {
	// 获取购物项中的商品信息
	commodityDao := dao.NewCommodityDao(cds.ctx)
	commodityIdList := lo.Map(cartItems, func(item *do.ShoppingCartItem, index int) int64 {
		return item.CommodityId
	})
	commodities, err := commodityDao.FindCommodities(commodityIdList)
	if err != nil {
		return errcode.Wrap("CartItemFillInCommodityInfoError", err)
	}
	if len(commodities) != len(cartItems) {
		logger.New(cds.ctx).Error("fillInCommodityError", "err", "商品信息不匹配", "commodityIdList", commodityIdList,
			"fetchedCommodities", commodities)
		return errcode.ErrCartItemParam
	}
	// 转换成以ID为Key的商品Map
	commodityMap := lo.SliceToMap(commodities, func(item *model.Commodity) (int64, *model.Commodity) {
		return item.ID, item
	})
	for _, cartItem := range cartItems {
		cartItem.CommodityName = commodityMap[cartItem.CommodityId].Name
		cartItem.CommodityImg = commodityMap[cartItem.CommodityId].CoverImg
		cartItem.CommoditySellingPrice = commodityMap[cartItem.CommodityId].SellingPrice
	}

	return nil
}

// DeleteUserCartItem 删除购物项
func (cds *CartDomainSvc) DeleteUserCartItem(cartItemId, userId int64) error {
	cartItemModel, _ := cds.cartDao.GetCartItemById(cartItemId)
	if cartItemModel == nil || cartItemModel.UserId != userId {
		logger.New(cds.ctx).Error("DataMatchError", "cartItem", cartItemModel, "cartItemId", cartItemId, "userId", userId)
		return errcode.ErrParams
	}
	err := cds.cartDao.DeleteAnCartItem(cartItemModel)
	if err != nil {
		err = errcode.Wrap("DeleteUserCartItemError", err)
	}

	return err
}

// GetCheckedCartItems 获取选中的购物项
func (cds *CartDomainSvc) GetCheckedCartItems(cartItemIds []int64, userId int64) ([]*do.ShoppingCartItem, error) {
	cartItemModels, err := cds.cartDao.FindCartItems(cartItemIds)
	if err != nil {
		err = errcode.Wrap("GetCheckedCartItemsError", err)
		return nil, err
	}
	// 确保购物项归属用户与请求用户一致
	userCartItemModels := lo.Filter(cartItemModels, func(item *model.ShoppingCartItem, index int) bool {
		return item.UserId == userId
	})
	if len(userCartItemModels) != len(cartItemIds) {
		return nil, errcode.ErrCartWrongUser
	}
	userCartItems := make([]*do.ShoppingCartItem, 0, len(userCartItemModels))
	err = util.CopyProperties(&userCartItems, &cartItemModels)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}
	// 填充购物项的商品信息
	cds.fillInCommodityInfo(userCartItems)
	return userCartItems, nil
}
