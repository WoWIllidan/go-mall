package dao

import (
	"context"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type OrderDao struct {
	ctx context.Context
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{ctx: ctx}
}

func (od *OrderDao) CreateOrder(tx *gorm.DB, order *do.Order) error {
	orderModel := new(model.Order)
	err := util.CopyProperties(orderModel, order)
	if err != nil {
		return errcode.ErrCoverData.WithCause(err)
	}

	err = tx.WithContext(od.ctx).Create(orderModel).Error
	if err != nil {
		return err
	}
	// 填充orderId
	order.ID = orderModel.ID
	order.Address.OrderId = orderModel.ID
	for _, item := range order.Items {
		item.OrderId = orderModel.ID
	}
	// 创建订单项
	err = od.createOrderItems(tx, order.Items)
	if err != nil {
		return err
	}
	// 创建订单地址
	err = od.createOrderAddress(tx, order.Address)
	return err
}

func (od *OrderDao) createOrderItems(tx *gorm.DB, orderItems []*do.OrderItem) error {
	orderItemModels := make([]*model.OrderItem, 0, len(orderItems))
	err := util.CopyProperties(&orderItemModels, &orderItems)
	if err != nil {
		return errcode.ErrCoverData.WithCause(err)
	}

	return tx.WithContext(od.ctx).Create(orderItemModels).Error
}

func (od *OrderDao) createOrderAddress(tx *gorm.DB, orderAddress *do.OrderAddress) error {
	orderAddressModel := new(model.OrderAddress)
	err := util.CopyProperties(orderAddressModel, orderAddress)
	if err != nil {
		return errcode.ErrCoverData.WithCause(err)
	}

	return tx.WithContext(od.ctx).Create(orderAddressModel).Error
}

// GetUserOrders 查询用户订单
func (od *OrderDao) GetUserOrders(userId int64, offset, returnSize int) (orders []*model.Order, totalRows int64, err error) {
	err = DB().WithContext(od.ctx).Where("user_id = ?", userId).
		Offset(offset).Limit(returnSize).
		Find(&orders).Error
	if err != nil {
		return
	}

	// 查询满足条件的记录数
	DB().WithContext(od.ctx).Model(model.Order{}).Where("user_id = ?", userId).Count(&totalRows)
	return
}

// GetMultiOrdersAddress 获取多个订单的地址, 返回以 orderId 为Key, 对应的订单地址为值的 Map
func (od *OrderDao) GetMultiOrdersAddress(orderIds []int64) (map[int64]*model.OrderAddress, error) {
	orderAddressList := make([]*model.OrderAddress, 0, len(orderIds))
	err := DB().WithContext(od.ctx).Where("order_id in (?)", orderIds).
		Find(&orderAddressList).Error
	if err != nil {
		return nil, err
	}
	// 转换成以orderId为Key的Map
	orderAddressMap := make(map[int64]*model.OrderAddress)
	orderAddressMap = lo.SliceToMap(orderAddressList, func(item *model.OrderAddress) (int64, *model.OrderAddress) {
		return item.OrderId, item
	})

	return orderAddressMap, nil
}

// GetMultiOrdersItems 获取多个订单对应的订单明细列表, 返回以 orderId 为Key, 对应的订单明细列表为值的 Map
func (od *OrderDao) GetMultiOrdersItems(orderIds []int64) (map[int64][]*model.OrderItem, error) {
	orderItems := make([]*model.OrderItem, 0)
	err := DB().WithContext(od.ctx).Where("order_id in (?)", orderIds).
		Find(&orderItems).Error
	if err != nil {
		return nil, err
	}

	// 转换成以orderId为Key, 订单明细列表为值的 Map
	orderItemsMap := make(map[int64][]*model.OrderItem)
	orderItemsMap = lo.GroupBy(orderItems, func(item *model.OrderItem) int64 {
		return item.OrderId
	})

	return orderItemsMap, nil
}

func (od *OrderDao) GetOrderByNo(orderNo string) (*model.Order, error) {
	order := new(model.Order)
	err := DB().WithContext(od.ctx).Where("order_no = ?", orderNo).
		Find(order).Error

	return order, err
}

func (od *OrderDao) GetOrderAddress(orderId int64) (*model.OrderAddress, error) {
	orderAddress := new(model.OrderAddress)
	err := DB().WithContext(od.ctx).Where("order_id = ?", orderId).
		Find(orderAddress).Error

	return orderAddress, err
}

func (od *OrderDao) GetOrderItems(orderId int64) ([]*model.OrderItem, error) {
	orderItems := make([]*model.OrderItem, 0)
	err := DB().WithContext(od.ctx).Where("order_id = ?", orderId).
		Find(&orderItems).Error

	return orderItems, err
}

// UpdateOrderStatus 更新订单状态
func (od *OrderDao) UpdateOrderStatus(orderId int64, status int) error {
	return DBMaster().WithContext(od.ctx).Model(model.Order{}).
		Where("id = ?", orderId).
		Update("order_status", status).Error
}
