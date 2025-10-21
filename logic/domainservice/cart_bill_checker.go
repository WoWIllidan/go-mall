package domainservice

import (
	"math"

	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/samber/lo"
)

type CartBillChecker struct {
	UserId        int64
	checkingItems []*do.ShoppingCartItem
	Coupon        struct { // 可用的优惠券
		CouponId      int64
		CouponName    string
		DiscountMoney int // 减免金额, 单位: 分
		Threshold     int // 使用门槛, 单位: 分, 设置成 1000 表示满10元可用
	}
	Discount struct { // 可用的满减活动券
		DiscountId    int64
		DiscountName  string
		DiscountMoney int
		Threshold     int
	}
	VipOffRate int // VIP的折扣  8 折  = 20% off

	handler cartBillCheckHandler
}

func NewCartBillChecker(items []*do.ShoppingCartItem, userId int64) *CartBillChecker {
	checker := new(CartBillChecker)
	checker.UserId = userId
	checker.checkingItems = items
	checker.handler = &checkerStarter{}
	// 通过责任链设置 要检查的各种优惠项
	checker.handler.SetNext(&couponChecker{}).
		SetNext(&discountChecker{}).
		SetNext(&vipChecker{})

	return checker
}

func (cbc *CartBillChecker) GetBill() (*do.CartBillInfo, error) {
	err := cbc.handler.RunChecker(cbc)
	if err != nil {
		return nil, errcode.Wrap("CartBillCheckerError", err)
	}
	// 计算商品使用减免前的总价
	originalTotalPrice := lo.Reduce(cbc.checkingItems, func(agg int, item *do.ShoppingCartItem, index int) int {
		return agg + item.CommoditySellingPrice
	}, 0)

	// VIP 能减免的金额
	vipDiscountMoney := int(math.Round(float64(originalTotalPrice) * float64(cbc.VipOffRate) / 100.0))

	totalPrice := originalTotalPrice - vipDiscountMoney
	if cbc.Coupon.Threshold != 0 && originalTotalPrice > cbc.Coupon.Threshold {
		// 满足优惠卷使用条件
		totalPrice -= cbc.Coupon.DiscountMoney
	}

	if cbc.Discount.Threshold != 0 && totalPrice > cbc.Discount.Threshold {
		// 满足使用满减券
		totalPrice -= cbc.Discount.DiscountMoney
	}
	billInfo := new(do.CartBillInfo)
	billInfo.Coupon = cbc.Coupon
	billInfo.Discount = cbc.Discount
	billInfo.VipDiscountMoney = vipDiscountMoney
	billInfo.TotalPrice = totalPrice
	billInfo.OriginalTotalPrice = originalTotalPrice

	return billInfo, nil
}

type cartBillCheckHandler interface {
	RunChecker(*CartBillChecker) error
	SetNext(cartBillCheckHandler) cartBillCheckHandler
	Check(*CartBillChecker) error
}

// 充当抽象类型，实现公共方法，抽象方法留给实现类自己实现
type cartCommonChecker struct {
	nextHandler cartBillCheckHandler
}

func (n *cartCommonChecker) SetNext(handler cartBillCheckHandler) cartBillCheckHandler {
	n.nextHandler = handler
	return handler
}

func (n *cartCommonChecker) RunChecker(billChecker *CartBillChecker) (err error) {
	// 由于go无继承的概念, 只能用组合，组合跟继承不一样，这里如果 cartCommonChecker 实现了 Check 方法，
	// 那么匿名嵌套它的具体处理类型，执行Execute的时候，调用的还是内部Next对象的Do方法
	// 调用不到外部类型的 Check 方法，所以 cartCommonChecker 不能实现 Check 方法
	if n.nextHandler != nil {
		if err = n.nextHandler.Check(billChecker); err != nil {
			err = errcode.Wrap("CartBillCheckerError", err)
			return
		}

		return n.nextHandler.RunChecker(billChecker)
	}

	return
}

// couponChecker 优惠卷 checker
type couponChecker struct {
	cartCommonChecker
}

func (cc *couponChecker) Check(cbc *CartBillChecker) error {
	// TODO: 查用户是否有可用优惠券
	// 这里是Mock逻辑
	// dao.GetUserCoupon(cbc.UserId)
	cbc.Coupon = struct {
		CouponId      int64
		CouponName    string
		DiscountMoney int
		Threshold     int
	}{
		CouponId:      1,
		DiscountMoney: 100, // 假设可用优惠券ID为1， 减免100
		Threshold:     100,
	}
	return nil
}

// discountChecker 折扣减免 checker
type discountChecker struct {
	cartCommonChecker
}

func (dc *discountChecker) Check(cbc *CartBillChecker) error {
	// TODO: 查用户是否有可用的减免活动
	// 这里是Mock逻辑
	// dao.GetApplicableDiscount(cbc.UserId)
	cbc.Discount = struct {
		DiscountId    int64
		DiscountName  string
		DiscountMoney int
		Threshold     int
	}{
		DiscountId:    1,
		DiscountMoney: 100, // 假设可用优惠券ID为1， 减免100
		Threshold:     1000,
	}
	return nil
}

type vipChecker struct {
	cartCommonChecker
}

func (vc *vipChecker) Check(cbc *CartBillChecker) error {
	// TODO: 判断用户是不是会员, 有没有会员折扣
	//if isVip(userId) {
	//  cbc.VipOffRate = 12 // 比如vip减免12%
	//  return nil
	//}
	cbc.VipOffRate = 0 // 不是vip不减免
	return nil
}

type checkerStarter struct {
	cartCommonChecker
}

func (s *checkerStarter) Check(cbc *CartBillChecker) (err error) {
	// 空Handler 这里什么也不做, 目的是让抽象类的 RunChecker 能启动调用链
	return
}
