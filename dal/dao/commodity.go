package dao

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/logic/do"
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

// InitCategoryData 初始化商品分类
func (cd *CommodityDao) InitCategoryData(categoryDos []*do.CommodityCategory) error {
	categoryModels := make([]*model.CommodityCategory, 0, len(categoryDos))
	util.CopyProperties(&categoryModels, &categoryDos)

	return cd.BulkCreateCommodityCategories(categoryModels)
}
