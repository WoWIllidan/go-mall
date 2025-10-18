package appservice

import (
	"context"

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
