package library

import "context"

type DemoLib struct {
	ctx context.Context
}

// NewDemoLib 创建时上层通过ctx 把 gin.Ctx传递过来
func NewDemoLib(ctx context.Context) *DemoLib {
	return &DemoLib{ctx: ctx}
}
