/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountDao_test.go
 Time    : 2018/9/11 17:42
 Author  : yanue
 Desc    : 

------------------------------- go ---------------------------------*/

package service

import (
	"github.com/jinzhu/gorm"
	"github.com/yanue/go-esport-common"
	"testing"
)

func init() {
}

func TestAccountService_Reg(t *testing.T) {
	// 表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "account_" + defaultTableName
	}
	// 全局禁用表名复数
	db.SingularTable(true)
	// 自动迁移模式
	common.Warn("aa")
}
