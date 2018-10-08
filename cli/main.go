/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : test.go
 Time    : 2018/9/13 17:15
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package main

import (
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/transport"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/proto"
)

func init() {
	client.DefaultClient.Init(
		client.Transport(
			transport.NewTransport(transport.Secure(true)),
		),
	)
}

var token string
var cl proto.AccountClient

func main() {
	// create a new service
	service := grpc.NewService()
	// Use the generated client stub
	cl = proto.NewAccountClient(common.MicroServiceNameAccount, service.Client())
	// parse command line flags
	service.Init()

	runClient()
}

func runClient() {
	sendVerifyCode()
	//loginByPhone()

	//loginByWechat()
	//loginByQQ()
	//loginByAccount()
}
