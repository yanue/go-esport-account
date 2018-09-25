/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : test.go
 Time    : 2018/9/13 17:15
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package main

import (
	"fmt"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/transport"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/proto"
	"github.com/yanue/go-esport-common/util"
	"golang.org/x/net/context"
)

func init() {
	client.DefaultClient.Init(
		client.Transport(
			transport.NewTransport(transport.Secure(true)),
		),
	)
}

func init() {

}

var token string

func main() {
	// create a new service
	service := grpc.NewService()
	// Use the generated client stub
	cl := proto.NewAccountClient(common.MicroServiceNameAccount, service.Client())
	// parse command line flags
	service.Init()

	login(cl)
	getUserInfo(cl)
	getAccountInfo(cl)
}

func login(cl proto.AccountClient) {
	pass, _ := util.Rsa.RsaEncryptPublic("111111")

	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		LoginType:  proto.ELoginType_PHONE,
		Account:    "yanue",
		Phone:      "18503002165",
		VerifyCode: "1234",
		Password:   string(pass),
		Os:         common.OS_IOS,
		DeviceId:   "98491eca-1111-4cea-beb4-859eb714296d",
	})

	if err != nil {
		fmt.Println("err:", err)
		return
	}
	token = rsp.Token
	fmt.Println("rsp:", common.OS_IOS, rsp)
}
func getUserInfo(cl proto.AccountClient) {
	//fmt.Println("ctx", ctx)
	// Make request
	rsp, err := cl.GetUserInfo(context.Background(), &proto.PInt32{Val: 1})

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("rsp:", rsp)
}

func getAccountInfo(cl proto.AccountClient) {
	// Set arbitrary headers in context
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"Authorization": token,
	})
	//fmt.Println("ctx", ctx)
	// Make request
	rsp, err := cl.GetAccountInfo(ctx, &proto.PNoParam{})

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("rsp:", rsp)
}
