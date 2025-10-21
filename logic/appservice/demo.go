package appservice

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/api/reply"
	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/cache"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/WoWBytePaladin/go-mall/logic/domainservice"
)

// 演示DEMO, 后期使用时删掉

type DemoAppSvc struct {
	ctx           context.Context
	demoDomainSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context) *DemoAppSvc {
	return &DemoAppSvc{
		ctx:           ctx,
		demoDomainSvc: domainservice.NewDemoDomainSvc(ctx),
	}
}

//func (das *DemoAppSvc)DoSomething() {
//	demo, err := das.demoDomainSvc.GetDemoEntity(id)
//	if err != nil {
//		logger.New(das.ctx).Error("DemoAppSvc DoSomething err", err)
//		return err
//	}
//	......
//}

// GetDemoIdentities 配置GORM时的演示方法, 显的有点脑残,
// 后面章节再解释怎么用ApplicationService 进行逻辑解耦
func (das *DemoAppSvc) GetDemoIdentities() ([]int64, error) {
	demos, err := das.demoDomainSvc.GetDemos()
	if err != nil {
		return nil, err
	}
	identities := make([]int64, 0, len(demos))

	for _, demo := range demos {
		identities = append(identities, demo.Id)
	}
	return identities, nil
}

func (das *DemoAppSvc) CreateDemoOrder(orderRequest *request.DemoOrderCreate) (*reply.DemoOrder, error) {
	demoOrderDo := new(do.DemoOrder)
	err := util.CopyProperties(demoOrderDo, orderRequest)
	if err != nil {
		errcode.Wrap("请求转换成demoOrderDo失败", err)
		return nil, err
	}
	demoOrderDo, err = das.demoDomainSvc.CreateDemoOrder(demoOrderDo)
	if err != nil {
		return nil, err
	}

	// TODO2 做一些其他的创建订单成功后的外围逻辑
	// 比如异步发送创建订单创建通知

	// 设置缓存和读取, 测试项目中缓存的使用, 没有其他任何意义
	cache.SetDemoOrder(das.ctx, demoOrderDo)
	cacheData, _ := cache.GetDemoOrder(das.ctx, demoOrderDo.OrderNo)
	logger.New(das.ctx).Info("redis data", "data", cacheData)

	replyDemoOrder := new(reply.DemoOrder)
	err = util.CopyProperties(replyDemoOrder, demoOrderDo)
	if err != nil {
		errcode.Wrap("demoOrderDo转换成replyDemoOrder失败", err)
		return nil, err
	}

	return replyDemoOrder, err
}

func (das *DemoAppSvc) InitCommodityCategoryData() error {
	cds := domainservice.NewCommodityDomainSvc(das.ctx)
	err := cds.InitCategoryData()
	return err
}
