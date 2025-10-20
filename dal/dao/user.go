package dao

import (
	"context"
	"errors"

	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"gorm.io/gorm"
)

type UserDao struct {
	ctx context.Context
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{ctx: ctx}
}

func (ud *UserDao) CreateUser(userInfo *do.UserBaseInfo, userPasswordHash string) (*model.User, error) {
	userModel := new(model.User)
	err := util.CopyProperties(userModel, userInfo)
	if err != nil {
		err = errcode.Wrap("UserDaoCreateUserError", err)
		return nil, err
	}
	userModel.Password = userPasswordHash

	err = DBMaster().WithContext(ud.ctx).Create(userModel).Error
	if err != nil {
		err = errcode.Wrap("UserDaoCreateUserError", err)
		return nil, err
	}
	return userModel, nil
}

func (ud *UserDao) FindUserByLoginNam(loginName string) (*model.User, error) {
	user := new(model.User)
	err := DB().Where(model.User{LoginName: loginName}).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return user, nil
}

func (ud *UserDao) FindUserById(userId int64) (*model.User, error) {
	user := new(model.User)
	err := DB().Where(model.User{ID: userId}).Find(&user).Error // Find 查找不到数据时不会返回 gorm.ErrRecordNotFound
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ud *UserDao) UpdateUser(user *model.User) error {
	err := DBMaster().Model(user).Updates(user).Error
	return err
}
