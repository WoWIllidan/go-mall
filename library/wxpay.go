package library

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/common/util/httptool"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/WoWBytePaladin/go-mall/resources"
)

type WxPayLib struct {
	ctx       context.Context
	payConfig WxtPayConfig
}

type WxtPayConfig struct {
	AppId           string
	MchId           string
	PrivateSerialNo string
	AesKey          string
	NotifyUrl       string
}

func NewWxPayLib(ctx context.Context, config WxtPayConfig) *WxPayLib {
	return &WxPayLib{
		ctx:       ctx,
		payConfig: config,
	}
}

const prePayApiUrl = "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"

type PrePayParam struct {
	AppId       string `json:"app_id"`
	MchId       string `json:"mch_id"`       // 商户号ID
	Description string `json:"description"`  // 商品描述
	OutTradeNo  string `json:"out_trade_no"` // 业务的订单号
	NotifyUrl   string `json:"notify_url"`   // 结果回调通知url
	Amount      struct {
		Total    int    `json:"total"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Payer struct {
		OpenId string `json:"open_id"`
	} `json:"payer"`
}

// WxPayInvokeInfo 前端用JSAPI调起支付的参数信息
// 微信支付文档: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_4.shtml
type WxPayInvokeInfo struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type WxPayNotifyResponse struct {
	CreateTime string              `json:"create_time"`
	Resource   WxPayNotifyResource `json:"resource"`
}
type WxPayNotifyResource struct {
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}

// WxPayNotifyResourceData 微信支付结果通知中解密后的resource数据
type WxPayNotifyResourceData struct {
	TransactionID string `json:"transaction_id"`
	Amount        struct {
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
	Mchid       string    `json:"mchid"`
	TradeState  string    `json:"trade_state"`
	BankType    string    `json:"bank_type"`
	SuccessTime time.Time `json:"success_time"`
	Payer       struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	OutTradeNo     string `json:"out_trade_no"`
	AppID          string `json:"AppID"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
	Attach         string `json:"attach"`
}

// CreateOrderPay 创建支付信息
// @param order *do.Order 业务的订单信息
// @param userOpenId string 用户的Openid
// @return payInvokeInfo *WxPayInvokeInfo 前端用于调起微信支付的参数
// @return err error 错误信息
func (wpl *WxPayLib) CreateOrderPay(order *do.Order, userOpenId string) (payInvokeInfo *WxPayInvokeInfo, err error) {
	// 创建预支付单
	// 微信支付文档：https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_1.shtml
	payDescription := fmt.Sprintf("GOMALL 商场购买%s 等商品", order.Items[0].CommodityName)
	prePayPram := &PrePayParam{
		AppId:       wpl.payConfig.AppId,
		MchId:       wpl.payConfig.MchId,
		Description: payDescription,
		OutTradeNo:  order.OrderNo,
		NotifyUrl:   wpl.payConfig.NotifyUrl,
	}
	prePayPram.Amount.Total = order.PayMoney
	prePayPram.Amount.Currency = "CNY"
	prePayPram.Payer.OpenId = userOpenId
	reqBody, _ := json.Marshal(prePayPram)
	token, err := wpl.getToken(http.MethodPost, string(reqBody), prePayApiUrl)
	if err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
		return
	}
	_, replyBody, err := httptool.Post(wpl.ctx, prePayApiUrl, reqBody, httptool.WithHeaders(map[string]string{
		"Authorization": "WECHATPAY2-SHA256-RSA2048 " + token,
	}))
	if err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
		return
	}
	prepayReply := struct {
		PrePayId string `json:"prepay_id"`
	}{}
	if err = json.Unmarshal(replyBody, &prepayReply); err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
		return
	}
	// 生成前端调起支付需要的参数
	payInvokeInfo, err = wpl.genPayInvokeInfo(prepayReply.PrePayId)
	if err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
	}
	return payInvokeInfo, nil
}

// ValidateNotifySignature 验证微信支付结果通知的签名
// 微信API文档: https://pay.weixin.qq.com/docs/merchant/development/interface-rules/signature-verification.html
// @param timeStamp 签名生成时间 从 HTTP 头 Wechatpay-Timestamp 获取
// @param nonce HTTP 头 Wechatpay-Nonce 中的应答随机串
// @param signature 微信支付的应答签名, 通过HTTP头Wechatpay-Signature 传递
// @param rawPost 原始请求体

func (wpl *WxPayLib) ValidateNotifySignature(timeStamp, nonce, signature, rawPost string) (verifyRes bool, err error) {
	signatureBytes, _ := base64.StdEncoding.DecodeString(signature)
	message := fmt.Sprintf("%s\n%s\n%s\n", timeStamp, nonce, rawPost)
	pemFileReader, err := resources.LoadResourceFile("wxp_pub.pem")
	if err != nil {
		err = errcode.Wrap("WxPayLibValidateCallBackSignatureError", err)
		return
	}
	publicKeyStr, _ := ioutil.ReadAll(pemFileReader)
	//pem解码
	block, _ := pem.Decode(publicKeyStr)
	//x509解码
	publicKeyInterface, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicKeyInterface.PublicKey.(*rsa.PublicKey)
	//验证数字签名
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, util.SHA256HashBytes(message), signatureBytes) //crypto.SHA1
	verifyRes = nil == err

	return verifyRes, err
}

// DecryptNotifyResourceData 解密微信支付通知中的resource数据
func (wpl *WxPayLib) DecryptNotifyResourceData(rawPost string) (notifyResourceData *WxPayNotifyResource, err error) {
	var notifyResponse WxPayNotifyResponse
	if err = json.Unmarshal([]byte(rawPost), &notifyResponse); err != nil {
		return notifyResourceData, errcode.Wrap("WxPayLibDecryptNotifyDataError", err)
	}

	aseKey := []byte(wpl.payConfig.AesKey)
	nonce := []byte(notifyResponse.Resource.Nonce)
	associatedData := []byte(notifyResponse.Resource.AssociatedData)
	ciphertext, err := base64.StdEncoding.DecodeString(notifyResponse.Resource.Ciphertext)
	if err != nil {
		return notifyResourceData, errcode.Wrap("WxPayLibDecryptNotifyDataError", err)
	}
	block, _ := aes.NewCipher(aseKey)
	aesGCM, _ := cipher.NewGCM(block)
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, associatedData)
	if err != nil {
		return notifyResourceData, err
	}
	err = json.Unmarshal(plaintext, &notifyResourceData)
	return
}

// genToken 生成微信支付的请求签名
// 文档：https://pay.weixin.qq.com/docs/merchant/development/interface-rules/signature-generation.html
func (wpl *WxPayLib) getToken(httMethod string, requestBody string, wxApiUrl string) (token string, err error) {

	urlPart, err := url.Parse(wxApiUrl)
	if err != nil {
		return token, err
	}
	canonicalUrl := urlPart.RequestURI()
	timestamp := time.Now().Unix()
	nonce := util.RandomString(32)
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n", httMethod, canonicalUrl, timestamp, nonce, requestBody)
	// 商户私有证书放在了 resources 目录下
	pemFileReader, err := resources.LoadResourceFile("wxpay.private.pem")
	if err != nil {
		return token, err
	}
	privateKey, err := ioutil.ReadAll(pemFileReader)
	if err != nil {
		return token, err
	}

	sha256MsgBytes := util.SHA256HashBytes(message)
	signBytes, err := util.RsaSignPKCS1v15(sha256MsgBytes, privateKey, crypto.SHA256)
	if err != nil {
		return token, err
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)

	token = fmt.Sprintf("mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		wpl.payConfig.MchId, nonce, timestamp, wpl.payConfig.PrivateSerialNo, sign)
	return token, nil
}

// 生成调起支付的参数
// 微信支付文档: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_4.shtml
func (wpl *WxPayLib) genPayInvokeInfo(prepayId string) (payInvokeInfo *WxPayInvokeInfo, err error) {
	payInvokeInfo = &WxPayInvokeInfo{
		AppId:     wpl.payConfig.AppId,
		TimeStamp: fmt.Sprintf("%v", time.Now().Unix()),
		NonceStr:  util.RandomString(32),
		Package:   "prepay_id=" + prepayId,
		SignType:  "RSA",
	}
	// 签名
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", payInvokeInfo.AppId, payInvokeInfo.TimeStamp, payInvokeInfo.NonceStr, payInvokeInfo.Package)
	pemFileReader, err := resources.LoadResourceFile("wxpay.private.pem")
	if err != nil {
		return
	}
	privateKey, err := ioutil.ReadAll(pemFileReader)
	if err != nil {
		return
	}
	sha245MsgBytes := util.SHA256HashBytes(message)
	signBytes, err := util.RsaSignPKCS1v15(sha245MsgBytes, privateKey, crypto.SHA256)
	if err != nil {
		return
	}
	payInvokeInfo.PaySign = base64.StdEncoding.EncodeToString(signBytes)

	return payInvokeInfo, nil
}
