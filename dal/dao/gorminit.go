package dao

import (
	"github.com/WoWBytePaladin/go-mall/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _DbMaster *gorm.DB
var _DbSlave *gorm.DB

// DB 返回只读实例
func DB() *gorm.DB {
	return _DbSlave
}

// DBMaster 返回主库实例
func DBMaster() *gorm.DB {
	return _DbMaster
}

func init() {
	//logger.New(context.TODO()).Info("database info", "db", config.Database)
	_DbMaster = initDB(config.Database.Master)
	_DbSlave = initDB(config.Database.Slave)
}

func getDialector(t, dsn string) gorm.Dialector {
	//switch t { 项目数据库需要加载多数据源时去掉注释
	//case "postgres":
	//	return postgres.Open(dsn)
	//default:
	//	return mysql.Open(dsn)
	//}
	return mysql.Open(dsn)
}

func initDB(option config.DbConnectOption) *gorm.DB {
	db, err := gorm.Open(
		getDialector(option.Type, option.DSN),
		&gorm.Config{
			Logger: NewGormLogger(),
		},
	)
	if err != nil {
		panic(err)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxOpenConns(option.MaxOpenConn)
	sqlDb.SetMaxIdleConns(option.MaxIdleConn)
	sqlDb.SetConnMaxLifetime(option.MaxLifeTime)
	if err = sqlDb.Ping(); err != nil {
		panic(err)
	}
	return db
}
