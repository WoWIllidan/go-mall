package domainservice

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/dao"
	"github.com/WoWBytePaladin/go-mall/logic/do"
)

// 演示DEMO, 后期使用时删掉

type DemoDomainSvc struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomainSvc(ctx context.Context) *DemoDomainSvc {
	return &DemoDomainSvc{
		ctx:     ctx,
		DemoDao: dao.NewDemoDao(ctx),
	}
}

// GetDemos 配置GORM时的演示方法
func (dds *DemoDomainSvc) GetDemos() ([]*do.DemoOrder, error) {
	demos, err := dds.DemoDao.GetAllDemos()
	if err != nil {
		err = errcode.Wrap("query entity error", err)
		return nil, err
	}

	demoOrders := make([]*do.DemoOrder, 0, len(demos))
	for _, demo := range demos {
		demoOrder := new(do.DemoOrder)
		util.CopyProperties(demoOrder, demo)
		demoOrders = append(demoOrders, demoOrder)
	}

	return demoOrders, nil
}

func (dds *DemoDomainSvc) CreateDemoOrder(demoOrder *do.DemoOrder) (*do.DemoOrder, error) {
	// 生成订单号  先随便写个
	demoOrder.OrderNo = "20240627596615375920904456"

	demoOrderModel, err := dds.DemoDao.CreateDemoOrder(demoOrder)
	if err != nil {
		err = errcode.Wrap("创建DemoOrder失败", err)
		return nil, err
	}

	// TODO1: 写订单快照
	// 这里一般要在事务里写订单商品快照表, 这个等后面做需求时再演示
	err = util.CopyProperties(demoOrder, demoOrderModel)
	// 返回领域对象
	return demoOrder, err
}
