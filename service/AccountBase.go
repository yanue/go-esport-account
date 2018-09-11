/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-基础初始化工作

------------------------------- go ---------------------------------*/

package service

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"go-esport-account/common"
	"go-esport-account/proto"
	"time"
)

var rpc *AccountRpc

/**
初始化服务
 */
func InitAccountService() {
	rpcService := micro.NewService(
		micro.Name("go.esport.account.srv"),
	)

	rpcService.Init()

	// 注册rpc
	go proto.RegisterAccountHandler(rpcService.Server(), InitRpc())

	err := rpcService.Run()
	if err != nil {
		log.Fatal("启动失败:%v", err.Error())
	}
}

/**
初始化服务
 */
func InitRpc() *AccountRpc {
	common.Info("InitAccountService")
	rpc = new(AccountRpc)
	rpc.account = new(AccountService)

	return rpc
}

/**
卸载服务
 */
func UnInitAccountService() {
	common.Notice("UnInitAccountService start")
	time.Sleep(1 * time.Minute)
	common.Notice("UnInitAccountService end")
}
