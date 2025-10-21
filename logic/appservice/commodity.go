package appservice

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/api/reply"
	"github.com/WoWBytePaladin/go-mall/common/app"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/logic/domainservice"
)

type CommodityAppSvc struct {
	ctx                context.Context
	commodityDomainSvc *domainservice.CommodityDomainSvc
}

func NewCommodityAppSvc(ctx context.Context) *CommodityAppSvc {
	return &CommodityAppSvc{
		ctx:                ctx,
		commodityDomainSvc: domainservice.NewCommodityDomainSvc(ctx),
	}
}

// GetCategoryHierarchy 获取按层级划分的商品分类
func (cas *CommodityAppSvc) GetCategoryHierarchy() []*reply.HierarchicCommodityCategory {
	categories := cas.commodityDomainSvc.GetHierarchicCategories()
	replyData := make([]*reply.HierarchicCommodityCategory, 0, len(categories))
	if len(categories) == 0 {
		return replyData
	}

	err := util.CopyProperties(&replyData, &categories)
	if err != nil {
		// 出错后 Depress Error, 不触发 500 Server Error
		// 这个错误判断经常用到, 使用一个统一的ErrorMsg, 方便在生产环境做日志监控 在生产环境上监控系统监控日志中的 ConvertDataError 关键字, 来实现主动告警
		logger.New(cas.ctx).Error(errcode.ErrCoverData.Msg(), "err", err)
		return replyData
	}
	return replyData
}

// GetSubCategories 按ParentId查询直接子分类
func (cas *CommodityAppSvc) GetSubCategories(parentId int64) []*reply.CommodityCategory {
	categories, err := cas.commodityDomainSvc.GetSubCategories(parentId)
	replyData := make([]*reply.CommodityCategory, 0, len(categories))
	log := logger.New(cas.ctx)
	if err != nil {
		// 有错误返回空列表, 不阻塞前端
		log.Error("CommodityAppSvcGetSubCategoriesError", "err", err)
		return replyData
	}

	err = util.CopyProperties(&replyData, &categories)
	if err != nil {
		log.Error(errcode.ErrCoverData.Msg(), "err", err)
		return replyData
	}
	return replyData
}

// GetCategoryCommodityList 获取分类商品列表
func (cas *CommodityAppSvc) GetCategoryCommodityList(categoryId int64, pagination *app.Pagination) ([]*reply.CommodityListElem, error) {
	categoryInfo := cas.commodityDomainSvc.GetCategoryInfo(categoryId)
	if categoryInfo == nil || categoryInfo.ID == 0 {
		return nil, errcode.ErrParams
	}
	commodityList, err := cas.commodityDomainSvc.GetCommodityListInCategory(categoryInfo, pagination)
	if err != nil {
		return nil, err
	}
	replyCommodityList := make([]*reply.CommodityListElem, 0, len(commodityList))
	err = util.CopyProperties(&replyCommodityList, &commodityList)
	if err != nil {
		// 包装成数据转换错误, 并设置错误原因
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return replyCommodityList, nil
}

// SearchCommodity 商品搜索
func (cas *CommodityAppSvc) SearchCommodity(keyword string, pagination *app.Pagination) ([]*reply.CommodityListElem, error) {
	commodityList, err := cas.commodityDomainSvc.SearchCommodity(keyword, pagination)
	if err != nil {
		return nil, err
	}
	replyCommodityList := make([]*reply.CommodityListElem, 0, len(commodityList))
	err = util.CopyProperties(&replyCommodityList, &commodityList)
	if err != nil {
		return nil, errcode.ErrCoverData.WithCause(err)
	}

	return replyCommodityList, nil
}

// CommodityInfo 商品详情
func (cas *CommodityAppSvc) CommodityInfo(commodityId int64) *reply.Commodity {
	commodityDO := cas.commodityDomainSvc.GetCommodityInfo(commodityId)
	if commodityDO == nil || commodityDO.ID == 0 {
		return nil
	}

	commodityInfo := new(reply.Commodity)
	util.CopyProperties(commodityInfo, commodityDO)
	return commodityInfo
}
