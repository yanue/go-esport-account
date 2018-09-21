/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-基础初始化工作

------------------------------- go ---------------------------------*/

package service

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/transport"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/proto"
	"time"
)

var rpc *AccountRpc
var cache *AccountCache
var orm *AccountOrm

const RedisPrefix = common.ServiceNameAccount + "_"

/**
初始化服务
 */
func InitAccountService(dbUser, dbAuth, dbAddr, dbName, redisAddr, redisPass string) {
	common.Logs.Info("InitAccountService")
	// db
	orm = new(AccountOrm)
	orm.initDb(dbUser, dbAuth, dbAddr, dbName)

	// redis
	cache = new(AccountCache)
	cache.initRedis(redisAddr, redisPass)

	// account service
	acct := new(AccountService)
	acct.cache = cache
	acct.orm = orm

	// rpc service
	rpc = new(AccountRpc)
	rpc.account = acct

	// 微服务
	rpcService := grpc.NewService(
		micro.Name(common.MicroServiceNameAccount),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		// setup a new transport with secure option
		micro.Transport(
			// create new transport
			transport.NewTransport(
				// set to automatically secure
				transport.Secure(true),
			),
		),
	)
	rpcService.Init()

	// 注册rpc
	go proto.RegisterAccountHandler(rpcService.Server(), rpc)

	// 运行
	err := rpcService.Run()
	if err != nil {
		panic("启动失败:" + err.Error())
	}
}

/**
卸载服务
 */
func UnInitAccountService() {
	// 关闭连接
	cache.redis.Close()
	orm.db.Close()

	common.Logs.Info("UnInitAccountService")
}
