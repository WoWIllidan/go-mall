package appservice

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/api/reply"
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
