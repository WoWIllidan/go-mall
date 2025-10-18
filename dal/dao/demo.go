package dao

import (
	"context"

	"github.com/WoWBytePaladin/go-mall/dal/model"
)

type DemoDao struct {
	ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
	return &DemoDao{ctx: ctx}
}

func (demo *DemoDao) GetAllDemos() (demos []*model.DemoOrder, err error) {

	err = DB().WithContext(demo.ctx).Find(&demos).Error
	if err != nil {
		return nil, err
	}

	return demos, err
}
