package controller

import (
	"errors"

	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/common/app"
	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/logger"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/logic/appservice"
	"github.com/gin-gonic/gin"
)

func RefreshUserToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		app.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	userSvc := appservice.NewUserAppSvc(c)
	token, err := userSvc.TokenRefresh(refreshToken)
	if err != nil {
		if errors.Is(err, errcode.ErrTooManyRequests) {
			// 客户端有并发刷新token
			app.NewResponse(c).Error(errcode.ErrTooManyRequests)
		} else {
			appErr := err.(*errcode.AppError)
			app.NewResponse(c).Error(appErr)
		}
		return
	}
	app.NewResponse(c).Success(token)
}

func RegisterUser(c *gin.Context) {
	userRequest := new(request.UserRegister)
	if err := c.ShouldBind(userRequest); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if !util.PasswordComplexityVerify(userRequest.Password) {
		// Validator验证通过后再应用 密码复杂度这样的特殊验证
		logger.New(c).Warn("RegisterUserError", "err", "密码复杂度不满足", "password", userRequest.Password)
		app.NewResponse(c).Error(errcode.ErrParams)
		return
	}
	// 注册用户
	userSvc := appservice.NewUserAppSvc(c)
	err := userSvc.UserRegister(userRequest)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNameOccupied) {
			app.NewResponse(c).Error(errcode.ErrUserNameOccupied)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).SuccessOk()
	return
}

func LoginUser(c *gin.Context) {
	loginRequest := new(request.UserLogin)
	if err := c.ShouldBindJSON(&loginRequest.Body); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if err := c.ShouldBindHeader(&loginRequest.Header); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	// 登录用户
	userSvc := appservice.NewUserAppSvc(c)
	token, err := userSvc.UserLogin(loginRequest)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNotRight) {
			app.NewResponse(c).Error(errcode.ErrUserNotRight)
		} else if errors.Is(err, errcode.ErrUserInvalid) {
			app.NewResponse(c).Error(errcode.ErrUserNotRight)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		logger.New(c).Error("LoginError", "err", err)
		return
	}

	app.NewResponse(c).Success(token)
	return
}

func LogoutUser(c *gin.Context) {
	userId := c.GetInt64("userId")
	platform := c.GetString("platform")
	userSvc := appservice.NewUserAppSvc(c)
	err := userSvc.UserLogout(userId, platform)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(c).SuccessOk()
}
