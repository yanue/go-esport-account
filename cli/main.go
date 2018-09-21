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

func main() {
	// create a new service
	service := grpc.NewService()
	// Use the generated client stub
	cl := proto.NewAccountClient(common.MicroServiceNameAccount, service.Client())
	// parse command line flags
	service.Init()

	getUserInfo(cl)
}

func login(cl proto.AccountClient) {
	pass, _ := util.Rsa.RsaEncryptPublic("111111")

	type name struct {
		name string
	}

	var n *name

	n = new(name)

	fmt.Println("", n.name)

	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		LoginType:  proto.ELoginType_ACCOUNT,
		Account:    "yanue",
		Phone:      "18503002165",
		VerifyCode: "1234",
		Password:   string(pass),
	})

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("rsp:", rsp)
}
func getUserInfo(cl proto.AccountClient) {
	// Set arbitrary headers in context
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"Authorization": "CgVIUzI1NhIFcHJvdG8.EAEYy_uS3QU.446d5f0332809cba1d490bb0fc7c32e704338cdd4ece1c083a96789843fd46aa",
	})
	//fmt.Println("ctx", ctx)
	// Make request
	rsp, err := cl.GetAccountInfo(ctx, &proto.PInt32{Val: 1})

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("rsp:", rsp)
}
