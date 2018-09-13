/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : test.go
 Time    : 2018/9/13 17:15
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/proto"
	"golang.org/x/net/context"
)

func main() {
	// create a new service
	service := micro.NewService()

	// parse command line flags
	service.Init()

	// Use the generated client stub
	cl := proto.NewAccountClient(common.MicroServiceNameAccount, service.Client())

	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		Phone: "18503002165",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp)
}
