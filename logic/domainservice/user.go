package domainservice

import (
	"context"
	"time"

	"github.com/WoWBytePaladin/go-mall/common/enum"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/cache"
	"github.com/WoWBytePaladin/go-mall/logic/do"
)

type UserDomainSvc struct {
	ctx context.Context
}

func NewUserDomainSvc(ctx context.Context) *UserDomainSvc {
	return &UserDomainSvc{ctx: ctx}
}

// GetUserBaseInfo 因为还没开发注册登录功能, 这里先Mock一个返回
func (us *UserDomainSvc) GetUserBaseInfo(userId int64) *do.UserBaseInfo {
	return &do.UserBaseInfo{
		ID:        12345678,
		Nickname:  "Kevin",
		LoginName: "kev@gomall.com",
		Verified:  1,
		Avatar:    "",
		Slogan:    "",
		IsBlocked: enum.UserBlockStateNormal,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GenAuthToken 生成AccessToken和RefreshToken
// 在缓存中会存储最新的Token 以及与Platform对应的 UserSession 同时会删除缓存中旧的Token-其中RefreshToken采用的是延迟删除
// **UserSession 在设置时会覆盖掉旧的Session信息
func (us *UserDomainSvc) GenAuthToken(userId int64, platform string, sessionId string) (*do.TokenInfo, error) {
	user := us.GetUserBaseInfo(userId)
	// 处理参数异常情况, 用户不存在、被删除、被禁用
	if user.ID == 0 || user.IsBlocked == enum.UserBlockStateBlocked {
		err := errcode.ErrUserInvalid
		return nil, err
	}
	userSession := new(do.SessionInfo)
	userSession.UserId = userId
	userSession.Platform = platform
	if sessionId == "" {
		// 为空是用户的登录行为, 重新生成sessionId
		sessionId = util.GenSessionId(userId)
	}
	userSession.SessionId = sessionId
	accessToken, refreshToken, err := util.GenUserAuthToken(userId)
	// 设置 userSession 缓存
	userSession.AccessToken = accessToken
	userSession.RefreshToken = refreshToken
	if err != nil {
		err = errcode.Wrap("Token生成失败", err)
		return nil, err
	}
	// 向缓存中设置AccessToken和RefreshToken的缓存
	err = cache.SetUserToken(us.ctx, userSession)
	if err != nil {
		errcode.Wrap("设置Token缓存时发生错误", err)
		return nil, err
	}
	err = cache.DelOldSessionTokens(us.ctx, userSession)
	if err != nil {
		errcode.Wrap("删除旧Token时发生错误", err)
		return nil, err
	}
	err = cache.SetUserSession(us.ctx, userSession)
	if err != nil {
		errcode.Wrap("设置Session缓存时发生错误", err)
		return nil, err
	}

	srvCreateTime := time.Now()
	tokenInfo := &do.TokenInfo{
		AccessToken:   userSession.AccessToken,
		RefreshToken:  userSession.RefreshToken,
		Duration:      int64(enum.AccessTokenDuration.Seconds()),
		SrvCreateTime: srvCreateTime,
	}

	return tokenInfo, nil

}

func (us *UserDomainSvc) RefreshToken(refreshToken string) (*do.TokenInfo, error) {
	log := logger.New(us.ctx)
	ok, err := cache.LockTokenRefresh(us.ctx, refreshToken)
	defer cache.UnlockTokenRefresh(us.ctx, refreshToken)
	if err != nil {
		err = errcode.Wrap("刷新Token时设置Redis锁发生错误", err)
		return nil, err
	}
	if !ok {
		err = errcode.ErrTooManyRequests
		return nil, err
	}
	tokenSession, err := cache.GetRefreshToken(us.ctx, refreshToken)
	if err != nil {
		log.Error("GetRefreshTokenCacheErr", "err", err)
		// 服务断发生错误一律提示客户端Token有问题
		// 生产环境可以做好监控日志中这个错误的监控
		err = errcode.ErrToken
		return nil, err
	}
	// refreshToken没有对应的缓存
	if tokenSession == nil || tokenSession.UserId == 0 {
		err = errcode.ErrToken
		return nil, err
	}
	userSession, err := cache.GetUserPlatformSession(us.ctx, tokenSession.UserId, tokenSession.Platform)
	if err != nil {
		log.Error("GetUserPlatformSessionErr", "err", err)
		err = errcode.ErrToken
		return nil, err
	}
	// 请求刷新的RefreshToken与UserSession中的不一致, 证明这个RefreshToken已经过时
	// RefreshToken被窃取或者前端页面刷Token不是串行的互斥操作都有可能造成这种情况
	if userSession.RefreshToken != refreshToken {
		// 记一条警告日志
		log.Warn("ExpiredRefreshToken", "requestToken", refreshToken, "newToken", userSession.RefreshToken, "userId", userSession.UserId)
		// 错误返回Token不正确, 或者更精细化的错误提示已在xxx登录如不是您本人操作请xxx
		err = errcode.ErrToken
		return nil, err
	}

	// 重新生成Token  因为不是用户主动登录所以sessionID与之前的保持一致
	tokenInfo, err := us.GenAuthToken(tokenSession.UserId, tokenSession.Platform, tokenSession.SessionId)
	if err != nil {
		err = errcode.Wrap("GenAuthTokenErr", err)
		return nil, err
	}
	return tokenInfo, nil
}

func (us *UserDomainSvc) VerifyAccessToken(accessToken string) (*do.TokenVerify, error) {
	tokenInfo, err := cache.GetAccessToken(us.ctx, accessToken)
	if err != nil {
		logger.New(us.ctx).Error("GetAccessTokenErr", "err", err)
		return nil, err
	}
	tokenVerify := new(do.TokenVerify)
	if tokenInfo != nil && tokenInfo.UserId != 0 {
		tokenVerify.UserId = tokenInfo.UserId
		tokenVerify.SessionId = tokenInfo.SessionId
		tokenVerify.Approved = true
	} else {
		tokenVerify.Approved = false
	}
	return tokenVerify, nil
}
