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
