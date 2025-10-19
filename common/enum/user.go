package enum

import "time"

const (
	UserBlockStateNormal  = 0
	UserBlockStateBlocked = 1
)

const AccessTokenDuration = 2 * time.Hour
const RefreshTokenDuration = 24 * time.Hour * 10
const OldRefreshTokenHoldingDuration = 6 * time.Hour // 刷新Token时老的RefreshToken保留的时间(用于发现refresh被窃取)
