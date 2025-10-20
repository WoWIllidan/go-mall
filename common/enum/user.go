package enum

import "time"

const (
	UserBlockStateNormal  = 0
	UserBlockStateBlocked = 1
)

const AccessTokenDuration = 2 * time.Hour
const RefreshTokenDuration = 24 * time.Hour * 10
const OldRefreshTokenHoldingDuration = 6 * time.Hour // 刷新Token时老的RefreshToken保留的时间(用于发现refresh被窃取)
const PasswordTokenDuration = 15 * time.Minute       // 重置密码的验证Token的有效期

const AddressIsNotUserDefault = 0
const AddressIsUserDefault = 1 // 用户收货地址状态--默认地址
