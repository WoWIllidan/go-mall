package dao

import (
	"context"
	"errors"

	"github.com/WoWBytePaladin/go-mall/common/enum"
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
	err := DB().WithContext(ud.ctx).Where(model.User{ID: userId}).Find(&user).Error // Find 查找不到数据时不会返回 gorm.ErrRecordNotFound
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ud *UserDao) UpdateUser(user *model.User) error {
	err := DBMaster().WithContext(ud.ctx).Model(user).Updates(user).Error
	return err
}

func (ud *UserDao) CreateUserAddress(address *do.UserAddressInfo) (*model.UserAddress, error) {
	addressModel := new(model.UserAddress)
	err := util.CopyProperties(addressModel, address)
	if err != nil {
		err = errcode.Wrap("UserDaoCreateUserAddressError", err)
		return nil, err
	}
	// 确定用户是否要更新默认地址
	var defaultAddress *model.UserAddress
	if addressModel.Default == enum.AddressIsUserDefault {
		defaultAddress, err = ud.GetUserDefaultAddress(addressModel.UserId)
		if err != nil {
			return nil, err
		}
	}

	if defaultAddress != nil && defaultAddress.ID != 0 { // 存在默认地址则把原默认地址更新为非默认,然后再写入新的地址信息
		err = DBMaster().Transaction(func(tx *gorm.DB) error {
			// 注意: Updates 方法会忽略结构体中字段的零值, 需要配合 Select 选择要更新成零值的字段名才能更新成功
			err := tx.WithContext(ud.ctx).Model(defaultAddress).Select("Default").
				Updates(model.UserAddress{Default: enum.AddressIsNotUserDefault}).Error
			if err != nil {
				return err
			}
			err = tx.WithContext(ud.ctx).Create(addressModel).Error
			return err
		})
	} else {
		err = DBMaster().WithContext(ud.ctx).Create(addressModel).Error
	}

	if err != nil {
		return nil, err
	}
	return addressModel, nil
}

func (ud *UserDao) UpdateUserAddress(address *do.UserAddressInfo) error {
	addressModel := new(model.UserAddress)
	err := util.CopyProperties(addressModel, address)
	if err != nil {
		err = errcode.Wrap("UserDaoUpdateUserAddressError", err)
		return err
	}

	// 确定用户是否要更新默认地址
	var defaultAddress *model.UserAddress
	if addressModel.Default == enum.AddressIsUserDefault {
		defaultAddress, err = ud.GetUserDefaultAddress(addressModel.UserId)
		if err != nil {
			return err
		}
	}
	if defaultAddress != nil && defaultAddress.ID != 0 && defaultAddress.ID != addressModel.ID {
		// 用户已有默认地址, 且默认地址信息不是要更新的这条地址信息, 先把取消原来默认地址的默认状态再对本条地址进行更新
		err = DBMaster().Transaction(func(tx *gorm.DB) error {
			err := tx.WithContext(ud.ctx).Model(defaultAddress).Select("Default").
				Updates(&model.UserAddress{Default: enum.AddressIsNotUserDefault}).Error
			if err != nil {
				return err
			}
			err = tx.WithContext(ud.ctx).Model(addressModel).Updates(addressModel).Error
			return err
		})
	} else {
		err = DBMaster().WithContext(ud.ctx).Model(addressModel).Updates(addressModel).Error
	}

	return err
}

func (ud *UserDao) FindUserAddresses(userId int64) ([]*model.UserAddress, error) {
	addresses := make([]*model.UserAddress, 0)
	err := DB().WithContext(ud.ctx).Where(model.UserAddress{UserId: userId}).
		Order("`default` DESC"). // 把默认地址排第一个
		Find(&addresses).Error
	return addresses, err
}

func (ud *UserDao) GetSingleAddress(addressId int64) (*model.UserAddress, error) {
	address := new(model.UserAddress)
	err := DB().WithContext(ud.ctx).Where(model.UserAddress{ID: addressId}).
		Find(&address).Error
	return address, err
}

func (ud *UserDao) GetUserDefaultAddress(userId int64) (*model.UserAddress, error) {
	address := new(model.UserAddress)
	err := DB().WithContext(ud.ctx).Where(model.UserAddress{UserId: userId, Default: enum.AddressIsUserDefault}).
		First(&address).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return address, nil
}

func (ud *UserDao) DeleteOneAddress(address *model.UserAddress) error {
	return DBMaster().WithContext(ud.ctx).Delete(address).Error
}
