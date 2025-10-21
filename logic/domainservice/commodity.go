package domainservice

import (
	"context"
	"encoding/json"
	"errors"

	"sort"

	"github.com/WoWBytePaladin/go-mall/common/app"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/dao"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/WoWBytePaladin/go-mall/resources"
)

type CommodityDomainSvc struct {
	ctx          context.Context
	commodityDao *dao.CommodityDao
}

func NewCommodityDomainSvc(ctx context.Context) *CommodityDomainSvc {
	return &CommodityDomainSvc{
		ctx:          ctx,
		commodityDao: dao.NewCommodityDao(ctx),
	}
}

// InitCategoryData 初始化分类信息测试数据
func (cds *CommodityDomainSvc) InitCategoryData() error {
	categories, err := cds.commodityDao.GetAllCategories()
	if err != nil {
		return errcode.Wrap("初始化商品分类错误", err)
	}
	if len(categories) > 1 {
		// 避免重复初始化
		return errcode.Wrap("重复初始化商品分类", errors.New("不能重复初始化商品分类"))
	}

	cateInitFileHandler, _ := resources.LoadResourceFile("category_init_data.json")

	categoryDos := make([]*do.CommodityCategory, 0, len(categories))
	decoder := json.NewDecoder(cateInitFileHandler)
	decoder.Decode(&categoryDos)

	err = cds.commodityDao.InitCategoryData(categoryDos)
	if err != nil {
		return errcode.Wrap("初始化商品分类错误", err)
	}

	return nil
}

// GetHierarchicCategories 返回按层级划分的商品分类
func (cds *CommodityDomainSvc) GetHierarchicCategories() []*do.HierarchicCommodityCategory {
	categoryModels, _ := cds.commodityDao.GetAllCategories()

	FlatCategories := make([]*do.HierarchicCommodityCategory, 0, len(categoryModels))
	util.CopyProperties(&FlatCategories, categoryModels)

	// 按照level asc, rank desc, id asc 进行排序
	sort.SliceStable(FlatCategories, func(i, j int) bool {
		if FlatCategories[i].Level != FlatCategories[j].Level {
			return FlatCategories[i].Level < FlatCategories[j].Level
		}
		if FlatCategories[i].Rank != FlatCategories[j].Rank {
			return FlatCategories[i].Rank > FlatCategories[j].Rank
		}
		return FlatCategories[i].ID < FlatCategories[j].ID
	})
	// 构造一个分类的临时Map, Key为 1, 2层级的分类ID, 值为分类信息(包含其子分类)
	categoryTempMap := make(map[int64]*do.HierarchicCommodityCategory)
	for _, category := range FlatCategories {
		if category.ParentId == 0 {
			categoryTempMap[category.ID] = category
		} else if category.ParentId != 0 && category.Level == 2 {
			categoryTempMap[category.ID] = category
			// 把分类添加到父分类的 subCategories 中
			categoryTempMap[category.ParentId].SubCategories =
				append(categoryTempMap[category.ParentId].SubCategories, category)
		} else if category.ParentId != 0 && category.Level == 3 {
			// 把分类添加到父分类的 subCategories 中
			categoryTempMap[category.ParentId].SubCategories =
				append(categoryTempMap[category.ParentId].SubCategories, category)
		}
	}
	//组装按层级划分的商品分类
	var hierarchicCategories []*do.HierarchicCommodityCategory
	for _, category := range FlatCategories {
		if category.ParentId != 0 {
			continue
		}
		category.SubCategories = categoryTempMap[category.ID].SubCategories
		for _, subCategory := range category.SubCategories {
			subCategory.SubCategories = categoryTempMap[subCategory.ID].SubCategories
		}
		hierarchicCategories = append(hierarchicCategories, category)
	}

	return hierarchicCategories
}

// GetSubCategories 获取ParentId对应的直接子分类
func (cds *CommodityDomainSvc) GetSubCategories(parentId int64) ([]*do.CommodityCategory, error) {
	categoryModels, err := cds.commodityDao.GetSubCategories(parentId)
	if err != nil {
		return nil, errcode.Wrap("GetSubCategoriesError", err)
	}
	categories := make([]*do.CommodityCategory, 0, len(categoryModels))
	util.CopyProperties(&categories, &categoryModels)

	return categories, nil
}

// GetCategoryInfo 获取分类ID对应的分类信息
func (cds *CommodityDomainSvc) GetCategoryInfo(categoryId int64) *do.CommodityCategory {
	categoryModel, err := cds.commodityDao.GetCategoryById(categoryId)
	if err != nil {
		logger.New(cds.ctx).Error("GetCategoryInfoError", err)
		return nil
	}

	categoryInfo := new(do.CommodityCategory)
	util.CopyProperties(&categoryInfo, &categoryModel)
	return categoryInfo
}

// InitCommodityData 初始化商品数据
func (cds *CommodityDomainSvc) InitCommodityData() error {
	commodity, err := cds.commodityDao.GetOneCommodity()
	if err != nil {
		return errcode.Wrap("初始化商品错误", err)
	}
	if commodity.ID > 0 {
		// 商品表里有数据, 打断流程, 避免重复初始化
		return errcode.Wrap("重复初始化商品", errors.New("不能重复初始化商品"))
	}

	initDataFileReader, err := resources.LoadResourceFile("commodity_init_data.json")
	if err != nil {
		return errcode.Wrap("初始化商品错误", err)
	}

	commodityDos := make([]*do.Commodity, 0)
	decoder := json.NewDecoder(initDataFileReader)
	decoder.Decode(&commodityDos)
	err = cds.commodityDao.InitCommodityData(commodityDos)
	if err != nil {
		return errcode.Wrap("初始化商品错误", err)
	}

	return nil
}

// GetCommodityListInCategory 获取分类下的商品列表
func (cds *CommodityDomainSvc) GetCommodityListInCategory(categoryInfo *do.CommodityCategory, pagination *app.Pagination) ([]*do.Commodity, error) {
	offset := pagination.Offset()
	size := pagination.GetPageSize()

	thirdLevelCategoryIds, err := cds.commodityDao.GetThirdLevelCategories(categoryInfo)
	if err != nil {
		return nil, errcode.Wrap("GetCommodityListInCategoryError", err)
	}
	commodityModelList, totalRows, err := cds.commodityDao.GetCommoditiesInCategory(thirdLevelCategoryIds, offset, size)
	if err != nil {
		return nil, errcode.Wrap("GetCommodityListInCategoryError", err)
	}
	pagination.SetTotalRows(int(totalRows))

	commodityList := make([]*do.Commodity, 0, len(commodityModelList))
	err = util.CopyProperties(&commodityList, &commodityModelList)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return commodityList, nil
}

// SearchCommodity 商品搜索
func (cds *CommodityDomainSvc) SearchCommodity(keyword string, pagination *app.Pagination) ([]*do.Commodity, error) {
	offset := pagination.Offset()
	size := pagination.GetPageSize()

	commodityModelList, totalRows, err := cds.commodityDao.FindCommodityWithNameKeyword(keyword, offset, size)
	if err != nil {
		return nil, errcode.Wrap("SearchCommodityError", err)
	}
	pagination.SetTotalRows(int(totalRows))

	commodityList := make([]*do.Commodity, 0, len(commodityModelList))
	err = util.CopyProperties(&commodityList, &commodityModelList)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return commodityList, nil
}

// GetCommodityInfo 获取商品详情
func (cds *CommodityDomainSvc) GetCommodityInfo(commodityId int64) *do.Commodity {
	commodityModel, err := cds.commodityDao.FindCommodityById(commodityId)
	log := logger.New(cds.ctx)
	if err != nil {
		log.Error("GetCommodityInfoError", "err", err)
		return nil
	}

	commodity := new(do.Commodity)
	err = util.CopyProperties(commodity, commodityModel)
	if err != nil {
		log.Error(errcode.ErrCoverData.Error(), "err", err)
		return nil
	}
	return commodity
}
