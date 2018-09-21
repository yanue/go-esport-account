/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-数据库操作

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/yanue/go-esport-common"
	"strings"
)

type AccountOrm struct {
	db *gorm.DB
}

// 初始化db
func (this *AccountOrm) initDb(dbUser, dbAuth, dbAddr, dbName string) {
	// user:password@tcp(addr)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	mysqlDsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbAuth, dbAddr, dbName)
	common.Logs.Info("mysqlDsn ", mysqlDsn)

	db, err := gorm.Open("mysql", mysqlDsn)
	if err != nil {
		panic("mysql连接失败:" + err.Error())
	}
	common.Logs.Info("db connected.")

	// 表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		// 替换t_
		return strings.Replace(defaultTableName, "t_", "", 1)
	}

	// 全局禁用表名复数
	db.SingularTable(true)
	// 自动迁移模式
	db.AutoMigrate(&TUser{})
	db.AutoMigrate(&TUserAuth{})
	db.AutoMigrate(&TAreaProvince{})
	db.AutoMigrate(&TAreaCity{})
	db.AutoMigrate(&TSchool{})
	db.AutoMigrate(&TClass{})

	orm.db = db
}

/**
 获取用户信息
 */
func (this *AccountOrm) GetUserInfo(uid int) (user *TUser, err error) {
	user = new(TUser)
	err = this.db.First(user, 10).Row().Scan(user)

	return user, err
}
