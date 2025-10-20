package domainservice

import (
	"context"
	"time"

	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/enum"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/cache"
	"github.com/WoWBytePaladin/go-mall/dal/dao"
	"github.com/WoWBytePaladin/go-mall/logic/do"
)

type UserDomainSvc struct {
	ctx     context.Context
	userDao *dao.UserDao
}

func NewUserDomainSvc(ctx context.Context) *UserDomainSvc {
	return &UserDomainSvc{
		ctx:     ctx,
		userDao: dao.NewUserDao(ctx),
	}
}

// GetUserBaseInfo UID查询用户信息
func (us *UserDomainSvc) GetUserBaseInfo(userId int64) *do.UserBaseInfo {
	user, err := us.userDao.FindUserById(userId)
	log := logger.New(us.ctx)
	if err != nil {
		log.Error("GetUserBaseInfoError", "err", err)
		return nil
	}
	userBaseInfo := new(do.UserBaseInfo)
	err = util.CopyProperties(userBaseInfo, user)
	if err != nil {
		log.Error("GetUserBaseInfoError", "err", err)
		return nil
	}
	return userBaseInfo
}

// UpdateUserBaseInfo 更新用户的基本信息
func (us *UserDomainSvc) UpdateUserBaseInfo(request *request.UserInfoUpdate, userId int64) error {
	user, err := us.userDao.FindUserById(userId)
	if err != nil {
		return err
	}

	user.Avatar = request.Avatar
	user.Nickname = request.Nickname
	user.Slogan = request.Slogan
	err = us.userDao.UpdateUser(user)
	return err
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
		tokenVerify.Platform = tokenInfo.Platform
		tokenVerify.Approved = true
	} else {
		tokenVerify.Approved = false
	}
	return tokenVerify, nil
}

func (us *UserDomainSvc) RegisterUser(userInfo *do.UserBaseInfo, plainPassword string) (*do.UserBaseInfo, error) {
	// 确定登录名可用
	existedUser, err := us.userDao.FindUserByLoginNam(userInfo.LoginName)
	if err != nil {
		return nil, errcode.Wrap("UserDomainSvcRegisterUserError", err)
	}
	if existedUser.LoginName != "" { // 用户名已经被占用
		return nil, errcode.ErrUserNameOccupied
	}
	passwordHash, err := util.BcryptPassword(plainPassword)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}
	userModel, err := us.userDao.CreateUser(userInfo, passwordHash)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}
	err = util.CopyProperties(userInfo, userModel)
	if err != nil {
		err = errcode.Wrap("UserDomainSvcRegisterUserError", err)
		return nil, err
	}

	return userInfo, nil
}

func (us *UserDomainSvc) LoginUser(loginName, plainPassword, platform string) (*do.TokenInfo, error) {
	existedUser, err := us.userDao.FindUserByLoginNam(loginName)
	if err != nil {
		return nil, errcode.Wrap("UserDomainSvcLoginUserError", err)
	}
	if existedUser.ID == 0 {
		return nil, errcode.ErrUserNotRight
	}
	if !util.BcryptCompare(existedUser.Password, plainPassword) {
		return nil, errcode.ErrUserNotRight
	}
	// 生成Token 和 Session
	tokenInfo, err := us.GenAuthToken(existedUser.ID, platform, "")
	return tokenInfo, err
}

func (us *UserDomainSvc) LogoutUser(userId int64, platform string) error {
	log := logger.New(us.ctx)
	userSession, err := cache.GetUserPlatformSession(us.ctx, userId, platform)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}
	// 删掉用户当前会话中的AccessToken和RefreshToken
	err = cache.DelAccessToken(us.ctx, userSession.AccessToken)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}
	err = cache.DelRefreshToken(us.ctx, userSession.RefreshToken)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}
	// 删掉用户在对应平台上的Session
	err = cache.DelUserSessionOnPlatform(us.ctx, userId, platform)
	if err != nil {
		log.Error("LogoutUserError", "err", err)
		return errcode.Wrap("UserDomainSvcLogoutUserError", err)
	}

	return nil
}

// ApplyForPasswordReset 申请重置密码
// @return passwordResetToken 重置密码时需要携带的Token信息，用于安全验证
// @return err 错误返回
func (us *UserDomainSvc) ApplyForPasswordReset(loginName string) (passwordResetToken, code string, err error) {
	user, err := us.userDao.FindUserByLoginNam(loginName)
	if err != nil {
		err = errcode.Wrap("ApplyForPasswordResetError", err)
		return
	}
	if user.ID == 0 {
		err = errcode.ErrUserNotRight
		return
	}
	token, err := util.GenPasswordResetToken(user.ID)
	code = util.RandNumStr(6)
	if err != nil {
		err = errcode.Wrap("ApplyForPasswordResetError", err)
		return
	}
	// 把token和验证码存入缓存
	err = cache.SetPasswordResetToken(us.ctx, user.ID, token, code)
	if err != nil {
		err = errcode.Wrap("ApplyForPasswordResetError", err)
		return
	}
	passwordResetToken = token
	return
}

func (us *UserDomainSvc) ResetPassword(resetToken, resetCode, newPlainPassword string) error {
	log := logger.New(us.ctx)
	userId, code, err := cache.GetPasswordResetToken(us.ctx, resetToken)
	if err != nil {
		log.Error("ResetPasswordError", "err", err)
		err = errcode.Wrap("ResetPasswordError", err)
		return err
	}
	// 确认Token正确且code码正确
	if userId == 0 || resetCode != code {
		return errcode.ErrParams
	}
	user, err := us.userDao.FindUserById(userId)
	if err != nil {
		return errcode.Wrap("ResetPasswordError", err)
	}
	// 找不到用户或者用户为封禁状态
	if user.ID == 0 || user.IsBlocked == enum.UserBlockStateBlocked {
		return errcode.ErrUserInvalid
	}
	newPass, err := util.BcryptPassword(newPlainPassword)
	if err != nil {
		return errcode.Wrap("ResetPasswordError", err)
	}
	// 更新密码
	user.Password = newPass
	err = us.userDao.UpdateUser(user)
	if err != nil {
		return errcode.Wrap("ResetPasswordError", err)
	}
	// 删掉用户所有已存的Session
	err = cache.DelUserSessions(us.ctx, userId)
	if err != nil {
		log.Error("ResetPasswordError", "err", err)
	}
	err = cache.DelPasswordResetToken(us.ctx, resetToken)
	if err != nil {
		// 删缓存失败, 不给客户端错误消息, 记日志发告警
		log.Error("ResetPasswordError", "err", err)
	}
	return nil
}

// AddUserAddress 新增用户收货地址
func (us *UserDomainSvc) AddUserAddress(addressInfo *do.UserAddressInfo) (*do.UserAddressInfo, error) {
	addressModel, err := us.userDao.CreateUserAddress(addressInfo)
	if err != nil {
		err = errcode.Wrap("AddUserAddressError", err)
		return nil, err
	}
	err = util.CopyProperties(addressInfo, addressModel)
	if err != nil {
		err = errcode.Wrap("AddUserAddressError", err)
		return nil, err
	}
	return addressInfo, nil
}

// GetUserAddresses 查询用户收货信息列表
func (us *UserDomainSvc) GetUserAddresses(userId int64) ([]*do.UserAddressInfo, error) {
	addresses, err := us.userDao.FindUserAddresses(userId)
	if err != nil {
		err = errcode.Wrap("GetUserAddressesError", err)
		return nil, err
	}
	userAddresses := make([]*do.UserAddressInfo, 0)
	if len(addresses) == 0 {
		return userAddresses, nil
	}
	err = util.CopyProperties(&userAddresses, &addresses)
	if err != nil {
		err = errcode.Wrap("AddUserAddressError", err)
		return nil, err
	}
	return userAddresses, nil
}

// GetUserSingleAddress 获取单个地址信息
func (us *UserDomainSvc) GetUserSingleAddress(userId, addressId int64) (*do.UserAddressInfo, error) {
	address, err := us.userDao.GetSingleAddress(addressId)
	if err != nil || address.UserId != userId {
		logger.New(us.ctx).Error("UserAddressNotMatchError", "err", err, "return data", address, "addressId", addressId, "userId", userId)
		return nil, errcode.ErrParams
	}
	userAddress := new(do.UserAddressInfo)
	err = util.CopyProperties(userAddress, address)
	if err != nil {
		return nil, errcode.Wrap("GetUserSingleAddressError", err)
	}
	return userAddress, nil
}

// ModifyUserAddress 更改用户的地址信息
func (us *UserDomainSvc) ModifyUserAddress(address *do.UserAddressInfo) error {
	addressModel, err := us.userDao.GetSingleAddress(address.ID)
	log := logger.New(us.ctx)
	if err != nil || address.UserId != addressModel.UserId {
		// 不匹配, 这种打一条日志, 监控系统按日志里的关键字做一下监控, 好发现问题
		log.Error("UserAddressNotMatchError", "err", err, "return data", address, "request data", address)
		return errcode.ErrParams
	}
	err = us.userDao.UpdateUserAddress(address)
	if err != nil {
		err = errcode.Wrap("UpdateUserAddressError", err)
	}
	return err
}

func (us *UserDomainSvc) DeleteOneUserAddress(userId, addressId int64) error {
	address, err := us.userDao.GetSingleAddress(addressId)
	if err != nil || address.UserId != userId {
		logger.New(us.ctx).Error("UserAddressNotMatchError", "err", err, "return data", address, "addressId", addressId, "userId", userId)
		return errcode.ErrParams
	}
	err = us.userDao.DeleteOneAddress(address)
	return err
}
