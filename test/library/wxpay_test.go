package library

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/library"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestWxPayLib_CreateOrderPay(t *testing.T) {
	defer gock.Off()
	openId := "QsudfrhgrDYDEEA1344EF"
	order := &do.Order{
		ID:          3,
		OrderNo:     "20240903374062590406950001",
		PayTransId:  "",
		PayType:     1,
		UserId:      1,
		BillMoney:   549900,
		PayMoney:    549700,
		PayState:    1,
		OrderStatus: 0,
		Address: &do.OrderAddress{
			OrderId:       0,
			UserName:      "张三",
			UserPhone:     "13512345679",
			ProvinceName:  "北京",
			CityName:      "北京",
			RegionName:    "朝阳区",
			DetailAddress: "XXX1号楼",
		},
		Items: []*do.OrderItem{
			{
				OrderId:               3,
				CommodityId:           2,
				CommodityName:         "Apple iPhone 11 (A2223)",
				CommodityImg:          "",
				CommoditySellingPrice: 599900,
				CommodityNum:          1,
			},
		},
		PaidAt:    time.Time{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	payConfig := library.WxtPayConfig{
		AppId:           "appId12345",
		MchId:           "mch12345",
		PrivateSerialNo: "567",
		AesKey:          "",
		NotifyUrl:       "",
	}
	payDescription := fmt.Sprintf("GOMALL 商场购买%s 等商品", order.Items[0].CommodityName)
	request := library.PrePayParam{
		AppId:       payConfig.AppId,
		MchId:       payConfig.MchId,
		Description: payDescription,
		OutTradeNo:  order.OrderNo,
		NotifyUrl:   payConfig.NotifyUrl,
		Amount: struct {
			Total    int    `json:"total"`
			Currency string `json:"currency"`
		}{
			Total:    order.PayMoney,
			Currency: "CNY",
		},
		Payer: struct {
			OpenId string `json:"open_id"`
		}{OpenId: openId},
	}

	gock.New("https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi").
		Post("").MatchType("json").
		JSON(request).
		Reply(200).
		JSON(map[string]string{"prepay_id": "wx26112221580621e9b071c00d9e093b0000"})

	wxPayLib := library.NewWxPayLib(context.TODO(), payConfig)
	var s *library.WxPayLib

	patchesOne := gomonkey.ApplyPrivateMethod(s, "getToken", func(_ *library.WxPayLib, httpMethod string, requestBody string, wxApiUrl string) (string, error) {
		token := fmt.Sprintf("mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
			payConfig.MchId, "abcddef", time.Now().Unix(), payConfig.PrivateSerialNo, "")
		return token, nil
	})

	patchsTwo := gomonkey.ApplyPrivateMethod(s, "genPayInvokeInfo", func(_ *library.WxPayLib) (*library.WxPayInvokeInfo, error) {
		payInfo := &library.WxPayInvokeInfo{
			AppId:     "123456",
			TimeStamp: fmt.Sprintf("%v", time.Now().Unix()),
			NonceStr:  "e61463f8efa94090b1f366cccfbbb444",
			Package:   "prepay_id=wx21201855730335ac86f8c43d1889123400",
			SignType:  "RSA",
			PaySign:   "oR9d8PuhnIc+YZ8cBHFCwfgpaK9gd7vaRvkYD7rthRAZ/X+QBhcCYL21N7cHCTUxbQ+EAt6Uy+lwSN22f5YZvI45MLko8Pfso0jm46v5hqcVwrk6uddkGuT+Cdvu4WBqDzaDjnNa5UK3GfE1Wfl2gHxIIY5lLdUgWFts17D4WuolLLkiFZV+JSHMvH7eaLdT9N5GBovBwu5yYKUR7skR8Fu+LozcSqQixnlEZUfyE55feLOQTUYzLmR9pNtPbPsu6WVhbNHMS3Ss2+AehHvz+n64GDmXxbX++IOBvm2olHu3PsOUGRwhudhVf7UcGcunXt8cqNjKNqZLhLw4jq/xDg==",
		}

		return payInfo, nil
	})
	defer patchesOne.Reset()
	defer patchsTwo.Reset()

	payInfo, err := wxPayLib.CreateOrderPay(order, openId)
	assert.Nil(t, err)
	logger.New(context.TODO()).Info("mock return", "payInfo", payInfo)
	assert.Equal(t, "e61463f8efa94090b1f366cccfbbb444", payInfo.NonceStr)
	if payInfo.PaySign == "" || payInfo.Package == "" {
		t.Fail()
	}
}
