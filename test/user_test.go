package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/WoWBytePaladin/go-mall/common/util"
	dao2 "github.com/WoWBytePaladin/go-mall/dal/dao"
	"github.com/WoWBytePaladin/go-mall/logic/do"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"regexp"
	"testing"
	"time"
)

var (
	mock sqlmock.Sqlmock
	err  error
	db   *sql.DB
)

// TestMain 是在当前package下，最先运行的一个函数，常用于测试基础组件的初始化
func TestMain(m *testing.M) {
	// 这里创建一个 sqlmock 的数据库连接 和 mock对象，mock对象管理 db 预期要执行的SQL

	// sqlmock 默认使用 sqlmock.QueryMatcherRegex 作为默认的SQL匹配器
	// 该匹配器使用mock.ExpectQuery 和 mock.ExpectExec 的参数作为正则表达式与真正执行的SQL语句进行匹配
	// 我们可以使用 regexp.QuoteMeta 把SQL转义成正则表达式 => mock.ExpectQuery(regexp.QuoteMeta("`SELECT * FROM `users`"))
	//
	// 如果想进行更严格的匹配, 可以让sqlmock 使用 sqlmock.QueryMatcherEqual 作为匹配器匹配器，该匹配器把mock.ExpectQuery
	// 和 mock.ExpectExec 的参数作为预期要执行的SQL语句跟真正要执行的SQL进行相等比较, 只有完全一样才会测试通过, 即使少个空格也不行
	// db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	// 把项目使用的DB连接换成sqlmock的DB连接
	dbMasterConn, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
		DefaultStringSize:         0,
	}))
	dbSlaveConn, _ := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
		DefaultStringSize:         0,
	}))
	dao2.SetDBMasterConn(dbMasterConn)
	dao2.SetDBSlaveConn(dbSlaveConn)

	// m.Run 是调用包下面各个Test函数的入口
	os.Exit(m.Run())
}

func TestUserDao_CreateUser(t *testing.T) {
	type fields struct {
		ctx context.Context
	}
	userInfo := &do.UserBaseInfo{
		Nickname:  "Slang",
		LoginName: "slang@go-mall.com",
		Verified:  0,
		Avatar:    "",
		Slogan:    "happy!",
		IsBlocked: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	passwordHash, _ := util.BcryptPassword("123456")
	userIsDel := 0

	ud := dao2.NewUserDao(context.TODO())
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(userInfo.Nickname, userInfo.LoginName, passwordHash, userInfo.Verified, userInfo.Avatar,
			userInfo.Slogan, userIsDel, userInfo.IsBlocked, userInfo.CreatedAt, userInfo.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	userObj, err := ud.CreateUser(userInfo, passwordHash)
	assert.Nil(t, err)
	assert.Equal(t, userInfo.LoginName, userObj.LoginName)
}
