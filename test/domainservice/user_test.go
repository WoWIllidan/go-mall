package domainservice

import (
	"context"
	"os"

	"github.com/WoWBytePaladin/go-mall/common/errcode"
	"github.com/WoWBytePaladin/go-mall/common/util"
	"github.com/WoWBytePaladin/go-mall/dal/dao"
	"github.com/WoWBytePaladin/go-mall/dal/model"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/WoWBytePaladin/go-mall/logic/domainservice"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// convey在TestMain下的入口
	SuppressConsoleStatistics()
	result := m.Run()
	// convey的结果打印
	PrintConsoleStatistics()
	os.Exit(result)
}

func TestPasswordComplexityVerify(t *testing.T) {
	Convey("Given a simple password", t, func() {
		password := "123456"
		Convey("When run it for password complexity checking", func() {
			result := util.PasswordComplexityVerify(password)
			Convey("Then the checking result should be false", func() {
				So(result, ShouldBeFalse)
			})
		})
	})

	Convey("Given a complex password", t, func() {
		password := "123@1~356Wrx"
		Convey("When run it for password complexity checking", func() {
			result := util.PasswordComplexityVerify(password)
			Convey("Then the checking result should be true", func() {
				So(result, ShouldBeTrue)
			})
		})
	})
}

func TestUserDomainSvc_RegisterUser(t *testing.T) {
	Convey("Given a user for RegisterUser of UserDomainSvc", t, func() {
		givenUser := &do.UserBaseInfo{
			Nickname:  "Kevin",
			LoginName: "kevin@go-mall.com",
			Verified:  0,
			Avatar:    "",
			Slogan:    "Keep tang ping",
			IsBlocked: 0,
			CreatedAt: time.Date(2025, 1, 31, 23, 28, 0, 0, time.Local),
			UpdatedAt: time.Date(2025, 1, 31, 23, 28, 0, 0, time.Local),
		}
		planPassword := "123@1~356Wrx"
		var s *dao.UserDao
		// 让UserDao的CreateUser返回Mock数据
		gomonkey.ApplyMethod(s, "CreateUser", func(_ *dao.UserDao, user *do.UserBaseInfo, password string) (*model.User, error) {
			passwordHash, _ := util.BcryptPassword(planPassword)
			userResult := &model.User{
				ID:        1,
				Nickname:  givenUser.Nickname,
				LoginName: givenUser.LoginName,
				Verified:  givenUser.Verified,
				Password:  passwordHash,
				Avatar:    givenUser.Avatar,
				Slogan:    givenUser.Slogan,
				CreatedAt: givenUser.CreatedAt,
				UpdatedAt: givenUser.UpdatedAt,
			}
			return userResult, nil
		})

		Convey("When the login name of user is not occupied", func() {
			gomonkey.ApplyMethod(s, "FindUserByLoginNam", func(_ *dao.UserDao, loginName string) (*model.User, error) {
				return new(model.User), nil
			})
			Convey("Then user should be created successfully", func() {
				user, err := domainservice.NewUserDomainSvc(context.TODO()).RegisterUser(givenUser, planPassword)
				So(err, ShouldBeNil)
				So(user.ID, ShouldEqual, 1)
				So(user, ShouldEqual, givenUser)
			})
		})

		Convey("When the login name of user has already been occupied by other users", func() {
			gomonkey.ApplyMethod(s, "FindUserByLoginNam", func(_ *dao.UserDao, loginName string) (*model.User, error) {
				return &model.User{LoginName: givenUser.LoginName}, nil
			})
			Convey("Then the user's registration should be unsuccessful", func() {
				user, err := domainservice.NewUserDomainSvc(context.TODO()).RegisterUser(givenUser, planPassword)
				So(user, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, errcode.ErrUserNameOccupied)
			})
		})
	})
}
