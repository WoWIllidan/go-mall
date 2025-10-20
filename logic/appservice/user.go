package appservice

import (
	"context"
	"errors"

	"github.com/WoWBytePaladin/go-mall/api/reply"
	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/WoWBytePaladin/go-mall/logic/domainservice"
)

type UserAppSvc struct {
	ctx           context.Context
	userDomainSvc *domainservice.UserDomainSvc
}

func NewUserAppSvc(ctx context.Context) *UserAppSvc {
	return &UserAppSvc{
		ctx:           ctx,
		userDomainSvc: domainservice.NewUserDomainSvc(ctx),
	}
}

func (us *UserAppSvc) GenToken() (*reply.TokenReply, error) {
	token, err := us.userDomainSvc.GenAuthToken(12345678, "h5", "")
	if err != nil {
		return nil, err
	}
	logger.New(us.ctx).Info("generate token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	util.CopyProperties(tokenReply, token)
	return tokenReply, err
}

func (us *UserAppSvc) TokenRefresh(refreshToken string) (*reply.TokenReply, error) {
	token, err := us.userDomainSvc.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	logger.New(us.ctx).Info("refresh token success", "tokenData", token)
	tokenReply := new(reply.TokenReply)
	util.CopyProperties(tokenReply, token)
	return tokenReply, err
}

func (us *UserAppSvc) UserRegister(userRegisterReq *request.UserRegister) error {
	userInfo := new(do.UserBaseInfo)
	util.CopyProperties(userInfo, userRegisterReq)

	// 调用领域服务注册用户
	_, err := us.userDomainSvc.RegisterUser(userInfo, userRegisterReq.Password)
	if errors.Is(err, errcode.ErrUserNameOccupied) {
		// 重名导致的注册不成功不需要额外处理
		return err
	}
	if err != nil && !errors.Is(err, errcode.ErrUserNameOccupied) {
		// TODO 发通知告知用户注册失败 ｜ 记录日志,监控告警,提示有用户注册失败发生
		return err
	}

	// TODO 写注册成功后的外围辅助逻辑, 比如注册成功后给用户发确认邮件|短信

	// TODO 如果产品逻辑是注册后帮用户登录, 那这里再掉登录的逻辑

	return nil
}

func (us *UserAppSvc) UserLogin(userLoginReq *request.UserLogin) (*reply.TokenReply, error) {
	tokenInfo, err := us.userDomainSvc.LoginUser(userLoginReq.Body.LoginName, userLoginReq.Body.Password, userLoginReq.Header.Platform)
	if err != nil {
		return nil, err
	}

	tokenReply := new(reply.TokenReply)
	util.CopyProperties(tokenReply, tokenInfo)

	// TODO 执行用户登录成功后发送消息通知之类的外围辅助型逻辑

	return tokenReply, nil
}

func (us *UserAppSvc) UserLogout(userId int64, platform string) error {
	err := us.userDomainSvc.LogoutUser(userId, platform)
	return err
}

// PasswordResetApply 申请重置密码
func (us *UserAppSvc) PasswordResetApply(request *request.PasswordResetApply) (*reply.PasswordResetApply, error) {
	passwordResetToken, code, err := us.userDomainSvc.ApplyForPasswordReset(request.LoginName)
	// TODO 把验证码通过邮件/短信发送给用户, 练习中就不实际去发送了, 记一条日志代替。
	logger.New(us.ctx).Info("PasswordResetApply", "token", passwordResetToken, "code", code)
	if err != nil {
		return nil, err
	}
	reply := new(reply.PasswordResetApply)
	reply.PasswordResetToken = passwordResetToken
	return reply, nil
}

// PasswordReset 重置密码
func (us *UserAppSvc) PasswordReset(request *request.PasswordReset) error {
	return us.userDomainSvc.ResetPassword(request.Token, request.Code, request.Password)
}

// UserInfo 用户信息
func (us *UserAppSvc) UserInfo(userId int64) *reply.UserInfoReply {
	userInfo := us.userDomainSvc.GetUserBaseInfo(userId)
	if userInfo == nil || userInfo.ID == 0 {
		return nil
	}
	infoReply := new(reply.UserInfoReply)
	util.CopyProperties(infoReply, userInfo)
	// 登录名是敏感信息, 做混淆处理
	infoReply.LoginName = util.MaskLoginName(infoReply.LoginName)
	return infoReply
}

// UserInfoUpdate 更新用户昵称、签名等信息
func (us *UserAppSvc) UserInfoUpdate(request *request.UserInfoUpdate, userId int64) error {
	return us.userDomainSvc.UpdateUserBaseInfo(request, userId)
}

// AddUserAddress 新增用户收货地址
func (us *UserAppSvc) AddUserAddress(request *request.UserAddress, userId int64) error {
	userAddressInfo := new(do.UserAddressInfo)
	err := util.CopyProperties(userAddressInfo, request)
	if err != nil {
		return errcode.Wrap("请求转换成领域对象失败", err)
	}
	userAddressInfo.UserId = userId
	newUserAddress, err := us.userDomainSvc.AddUserAddress(userAddressInfo)
	if err != nil {
		logger.New(us.ctx).Error("添加用户收货地址失败", "err", err, "return data", newUserAddress)
	}
	return err
}

// GetUserAddresses 查询用户所有收货地址信息
func (us *UserAppSvc) GetUserAddresses(userId int64) ([]*reply.UserAddress, error) {
	userAddresses := make([]*reply.UserAddress, 0)
	addresses, err := us.userDomainSvc.GetUserAddresses(userId)
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 { // 没有数据, 返回userAddressesReply 而不是nil, 避免格式化时data字段值为null
		return userAddresses, nil
	}
	err = util.CopyProperties(&userAddresses, &addresses)
	if err != nil {
		errcode.Wrap("GetUserAddresses转换响应数据时发生错误", err)
		return nil, err
	}
	for _, address := range userAddresses {
		// 用户姓名和手机号脱敏
		address.MaskedUserName = util.MaskRealName(address.UserName)
		address.MaskedUserPhone = util.MaskPhone(address.UserPhone)
	}
	return userAddresses, nil
}

// GetUserSingleAddress 获取单个地址信息
func (us *UserAppSvc) GetUserSingleAddress(userId, addressId int64) (*reply.UserAddress, error) {
	addressInfo, err := us.userDomainSvc.GetUserSingleAddress(userId, addressId)
	if err != nil {
		return nil, err
	}
	userAddress := new(reply.UserAddress)
	util.CopyProperties(userAddress, addressInfo)
	userAddress.MaskedUserName = util.MaskRealName(userAddress.UserName)
	userAddress.MaskedUserPhone = util.MaskPhone(userAddress.UserPhone)

	return userAddress, nil
}

// ModifyUserAddress 更新用户的某个收货地址信息
func (us *UserAppSvc) ModifyUserAddress(requestData *request.UserAddress, userId, addressId int64) error {
	userAddressInfo := new(do.UserAddressInfo)
	err := util.CopyProperties(userAddressInfo, requestData)
	if err != nil {
		return errcode.Wrap("请求转换成领域对象失败", err)
	}
	userAddressInfo.UserId = userId
	userAddressInfo.ID = addressId
	err = us.userDomainSvc.ModifyUserAddress(userAddressInfo)
	return err
}

// DeleteOneUserAddress 删除地址
func (us *UserAppSvc) DeleteOneUserAddress(userId, addressId int64) error {
	return us.userDomainSvc.DeleteOneUserAddress(userId, addressId)
}
