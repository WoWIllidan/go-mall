package request

import "time"

type OrderCreate struct {
	CartItemIdList []int64 `json:"cart_item_id_list" binding:"required"`
	UserAddressId  int64   `json:"user_address_id" binding:"required"`
}

// OrderPayCreate 订单发起支付请求
type OrderPayCreate struct {
	OrderNo string `json:"order_no" binding:"required"`
	PayType int    `json:"pay_type" binding:"required,oneof= 1 2"`
}

// WxPayNotifyRequest 微信支付回调通知请求
// https://pay.weixin.qq.com/docs/merchant/apis/jsapi-payment/payment-notice.html
type WxPayNotifyRequest struct {
	Header struct {
		Timestamp string `json:"Wechatpay-Timestamp"`
		Nonce     string `json:"Wechatpay-Nonce"`
		Signature string `json:"Wechatpay-Signature""`
	}
	Body struct {
		ID           string    `json:"id"`
		CreateTime   time.Time `json:"create_time"`
		ResourceType string    `json:"resource_type"`
		EventType    string    `json:"event_type"`
		Summary      string    `json:"summary"`
		Resource     struct {
			OriginalType   string `json:"original_type"`
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			AssociatedData string `json:"associated_data"`
			Nonce          string `json:"nonce"`
		} `json:"resource"`
	}
}
