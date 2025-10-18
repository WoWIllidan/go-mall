package domainservice

import (
	"context"
)

// 演示DEMO, 后期使用时删掉

type DemoDomainSvc struct {
	ctx context.Context
}

func NewDemoDomainSvc(ctx context.Context) *DemoDomainSvc {
	return &DemoDomainSvc{ctx: ctx}
}

//func GetDemoEntity(id int) (*DemoEntity, error){
//	entity, err := dao.GetDemoById(id)
//	// 转换成领域对象
//	copy(entity, domainEntity)
//	if err != nil {
//		err = errcode.Wrap("query entity error", err)
//		return nil, err
//	}
//	return domainEntity, nil
//}
