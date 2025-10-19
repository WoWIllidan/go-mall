package library

import (
	"context"
	"encoding/json"

	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util/httptool"
)

type DemoLib struct {
	ctx context.Context
}

// NewDemoLib 创建时上层通过ctx 把 gin.Ctx传递过来
func NewDemoLib(ctx context.Context) *DemoLib {
	return &DemoLib{ctx: ctx}
}

type OrderCreateResult struct {
	UserId    int64  `json:"user_id"`
	BillMoney int64  `json:"bill_money"`
	OrderNo   string `json:"order_no"`
	State     int8   `json:"state"`
	PaidAt    string `json:"paid_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 用http调自己项目里的POST接口, 演示用, 实际使用时不要这么干

func (lib *DemoLib) TestPostCreateOrder() (*OrderCreateResult, error) {
	data := &request.DemoOrderCreate{
		UserId:       12345,
		BillMoney:    20,
		OrderGoodsId: 1111110,
	}
	jsonReq, _ := json.Marshal(data)
	httCode, respBody, err := httptool.Post(lib.ctx, "http://localhost:8080/building/create-demo-order", jsonReq)
	logger.New(lib.ctx).Info("create-demo-order api response ", "code", httCode, "data", respBody, "err", err)

	if err != nil {
		return nil, err
	}

	reply := &struct {
		Code int `json:"code"`
		Data *OrderCreateResult
	}{}
	json.Unmarshal(respBody, reply)
	return reply.Data, nil
}
