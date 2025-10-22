package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type CommodityDao struct {
	ctx context.Context
}

func NewCommodityDao(ctx context.Context) *CommodityDao {
	return &CommodityDao{ctx: ctx}
}

func (cd *CommodityDao) BulkCreateCommodityCategories(categories []*model.CommodityCategory) error {
	return DBMaster().WithContext(cd.ctx).Create(categories).Error
}

func (cd *CommodityDao) BulkCreateCommodities(commodities []*model.Commodity) error {
	return DBMaster().WithContext(cd.ctx).Create(commodities).Error
}

func (cd *CommodityDao) GetAllCategories() ([]*model.CommodityCategory, error) {
	categories := make([]*model.CommodityCategory, 0)
	err := DB().WithContext(cd.ctx).Find(&categories).Error
	return categories, err
}

func (cd *CommodityDao) GetSubCategories(parentId int64) ([]*model.CommodityCategory, error) {
	categories := make([]*model.CommodityCategory, 0)
	err := DB().WithContext(cd.ctx).
		Where("parent_id = ?", parentId).
		Order("rank DESC").Find(&categories).Error
	return categories, err
}

// GetCategoryById 获取Id对应的分类信息
func (cd *CommodityDao) GetCategoryById(categoryId int64) (*model.CommodityCategory, error) {
	category := new(model.CommodityCategory)
	err := DB().WithContext(cd.ctx).Where("id = ?", categoryId).Find(category).Error
	return category, err
}

// GetSubCategoryIdList 查询分类的子分类ID
func (cd *CommodityDao) getSubCategoryIdList(parentCategoryIds []int64) (categoryIds []int64, err error) {
	err = DB().WithContext(cd.ctx).Model(model.CommodityCategory{}).
		Where("parent_id IN (?)", parentCategoryIds).
		Order("rank DESC").Pluck("id", &categoryIds).Error

	return
}

// InitCategoryData 初始化商品分类
func (cd *CommodityDao) InitCategoryData(categoryDos []*do.CommodityCategory) error {
	categoryModels := make([]*model.CommodityCategory, 0, len(categoryDos))
	util.CopyProperties(&categoryModels, &categoryDos)

	return cd.BulkCreateCommodityCategories(categoryModels)
}

// GetOneCommodity 无查询条件, 返回一条数据
func (cd *CommodityDao) GetOneCommodity() (*model.Commodity, error) {
	commodity := new(model.Commodity)
	err := DB().WithContext(cd.ctx).Find(commodity).Error
	return commodity, err
}

// InitCommodityData 初始化商品数据
func (cd *CommodityDao) InitCommodityData(commodities []*do.Commodity) error {
	commodityModels := make([]*model.Commodity, 0, len(commodities))
	util.CopyProperties(&commodityModels, &commodities)

	return cd.BulkCreateCommodities(commodityModels)
}

// GetThirdLevelCategories 查找分类下的所有三级分类ID
func (cd *CommodityDao) GetThirdLevelCategories(categoryInfo *do.CommodityCategory) (categoryIds []int64, err error) {
	if categoryInfo.Level == 3 {
		return []int64{categoryInfo.ID}, nil
	} else if categoryInfo.Level == 2 {
		categoryIds, err = cd.getSubCategoryIdList([]int64{categoryInfo.ID})
		return
	} else if categoryInfo.Level == 1 {
		var secondCategoryId []int64
		secondCategoryId, err = cd.getSubCategoryIdList([]int64{categoryInfo.ID})
		if err != nil {
			return
		}
		categoryIds, err = cd.getSubCategoryIdList(secondCategoryId)
		return
	}
	return
}

// GetCommoditiesInCategory 查询分类下的商品列表
func (cd *CommodityDao) GetCommoditiesInCategory(categoryIds []int64, offset, returnSize int) (commodityList []*model.Commodity, totalRows int64, err error) {
	// 查询满足条件的商品
	err = DB().WithContext(cd.ctx).Omit("detail_content"). // 忽略掉商品详情的内容
								Where("category_id IN (?)", categoryIds).
								Offset(offset).Limit(returnSize).
								Find(&commodityList).Error
	// 查询满足条件的记录数
	DB().WithContext(cd.ctx).Model(model.Commodity{}).Where("category_id IN (?)", categoryIds).Count(&totalRows)

	return
}

// FindCommodityWithNameKeyword 按名称LIKE查询商品列表
func (cd *CommodityDao) FindCommodityWithNameKeyword(keyword string, offset, returnSize int) (commodityList []*model.Commodity, totalRows int64, err error) {
	err = DB().WithContext(cd.ctx).Omit("detail_content").
		Where("name LIKE ?", "%"+keyword+"%").
		Offset(offset).Limit(returnSize).
		Find(&commodityList).Error
	DB().WithContext(cd.ctx).Model(model.Commodity{}).Where("name LIKE ?", "%"+keyword+"%").Count(&totalRows)

	return
}

// FindCommodityById 通过ID查商品信息
func (cd *CommodityDao) FindCommodityById(commodityId int64) (*model.Commodity, error) {
	commodity := new(model.Commodity)
	err := DB().WithContext(cd.ctx).Where("id = ?", commodityId).Find(commodity).Error
	return commodity, err
}

// FindCommodities 查询主键 id IN commodityIdList 的 商品
func (cd *CommodityDao) FindCommodities(commodityIdList []int64) ([]*model.Commodity, error) {
	commodities := make([]*model.Commodity, 0)
	err := DB().WithContext(cd.ctx).Find(&commodities, commodityIdList).Error
	return commodities, err
}

// ReduceStuckInOrderCreate 创建订单后商品减库存
func (cd *CommodityDao) ReduceStuckInOrderCreate(tx *gorm.DB, orderItems []*do.OrderItem) error {
	for _, orderItem := range orderItems {
		commodity := new(model.Commodity)
		// SELECT FOR UPDATE 当前读
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(cd.ctx).
			Find(commodity, orderItem.CommodityId)
		newStock := commodity.StockNum - orderItem.CommodityNum
		if newStock < 0 {
			return errcode.ErrCommodityStockOut.WithCause(errors.New("商品缺少库存, 商品ID:" + strconv.FormatInt(commodity.ID, 10)))
		}
		commodity.StockNum = newStock
		// https://gorm.io/docs/update.html#Update-single-column
		err := tx.WithContext(cd.ctx).Model(commodity).Update("stock_num", newStock).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// RecoverOrderCommodityStuck  用户取消订单后商品减库存
func (cd *CommodityDao) RecoverOrderCommodityStuck(orderItems []*do.OrderItem) error {
	err := DBMaster().Transaction(func(tx *gorm.DB) error {
		for _, orderItem := range orderItems {
			commodity := new(model.Commodity)
			tx.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(cd.ctx).Find(commodity, orderItem.CommodityId)
			if commodity.ID == 0 {
				return errcode.ErrNotFound.WithCause(errors.New(fmt.Sprintf("商品未找到, ID: %d", orderItem.CommodityId)))
			}
			newStock := commodity.StockNum + orderItem.CommodityNum
			err := tx.WithContext(cd.ctx).Model(commodity).Update("stock_num", newStock).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
