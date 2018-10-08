/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-redis操作

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/sms"
	"testing"
)

var acct *AccountService

func init() {
	// 读取配置信息
	dbUser := "root"
	dbAuth := "blemobi3721"
	dbAddr := "192.168.5.202"
	dbName := "esport-account"
	redisAddr := "192.168.5.201:6379"
	redisPass := ""

	common.Logs.Info("InitAccountService")
	// db
	orm = new(AccountOrm)
	orm.initDb(dbUser, dbAuth, dbAddr, dbName)

	// redis
	cache = new(AccountCache)
	cache.initRedis(redisAddr, redisPass)

	// account service
	acct = new(AccountService)
	acct.cache = cache
	acct.orm = orm

	// rpc service
	rpc = new(AccountRpc)
	rpc.account = acct

	accessKeyId := "LTAIiCR1y6RAa2IC"
	accessKeySecret := "pit2WJpgdhSwOEhzr42EdlXMuTdhpn"
	signName := "智享协同"
	// sms util
	smsUtil = sms.NewSms(accessKeyId, accessKeySecret, signName, cache.redis)
}

func TestAccountCache_GetUserInfo(t *testing.T) {
	//info, _ := acct.GetUserInfo(1)
	//fmt.Println("info", info)
	//fields, err := util.Struct.RangeToMap(info)
	//fmt.Println("fields, err", fields, err)
	//cache.redis.HMSet(cache.key.HUserInfo(1), fields)
	res, err := cache.GetUserInfo(1)
	fmt.Println("res, err", res, err)
}
