package dao

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/logic/do"
)

type CartDao struct {
	ctx context.Context
}

func NewCartDao(ctx context.Context) *CartDao {
	return &CartDao{ctx: ctx}
}

func (cd *CartDao) GetUserCartItemWithCommodityId(userId, commodityId int64) (*model.ShoppingCartItem, error) {
	var cartItem *model.ShoppingCartItem
	err := DB().WithContext(cd.ctx).Where(
		model.ShoppingCartItem{UserId: userId, CommodityId: commodityId},
		"UserId", "CommodityId"). // 保证Struct中的UserId, CommodityId为零值时仍用他们构建查询条件
		Find(&cartItem).Error
	return cartItem, err
}

func (cd *CartDao) GetCartItemById(id int64) (*model.ShoppingCartItem, error) {
	cartItem := new(model.ShoppingCartItem)
	err := DB().WithContext(cd.ctx).Where(model.ShoppingCartItem{CartItemId: id}, "CartItemId").Find(&cartItem).Error
	return cartItem, err
}

// FindCartItems 获取多个ID指定的购物项
func (cd *CartDao) FindCartItems(cartItemIdList []int64) ([]*model.ShoppingCartItem, error) {
	items := make([]*model.ShoppingCartItem, 0)
	// 查询主键 id IN cartItemIdList 的购物项
	err := DB().WithContext(cd.ctx).Find(&items, cartItemIdList).Error
	return items, err
}

func (cd *CartDao) AddCartItem(cartItem *do.ShoppingCartItem) error {
	cartItemModel := new(model.ShoppingCartItem)
	util.CopyProperties(cartItemModel, cartItem)

	err := DBMaster().WithContext(cd.ctx).Create(cartItemModel).Error
	return err
}

func (cd *CartDao) UpdateCartItem(cartItem *model.ShoppingCartItem) error {
	return DBMaster().WithContext(cd.ctx).Model(cartItem).Updates(cartItem).Error
}

func (cd *CartDao) GetUserCartItems(userId int64) ([]*model.ShoppingCartItem, error) {
	cartItems := make([]*model.ShoppingCartItem, 0)
	err := DB().WithContext(cd.ctx).Where(model.ShoppingCartItem{UserId: userId}, "UserId").
		Find(&cartItems).Error

	return cartItems, err
}

func (cd *CartDao) DeleteAnCartItem(cartItem *model.ShoppingCartItem) error {
	return DBMaster().WithContext(cd.ctx).Delete(cartItem).Error
}
