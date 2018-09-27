/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : api.go
 Time    : 2018/9/27 10:47
 Author  : yanue
 
 - 
 
------------------------------- go ---------------------------------*/

package main

import (
	"fmt"
	"github.com/micro/go-micro/metadata"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/proto"
	"github.com/yanue/go-esport-common/util"
	"golang.org/x/net/context"
)

func sendVerifyCode() {
	//fmt.Println("ctx", ctx)
	// Make request
	rsp, err := cl.SendSmsVerifyCode(context.Background(), &proto.PSmsData{
		Phone:    "18503002165",
		Imei:     "12323331231",
		CodeType: proto.PSmsData_quick_login,
	})

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("rsp:", rsp)
}

func login() {
	pass, _ := util.Rsa.RsaEncryptPublic("111111")

	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		LoginType:  proto.ELoginType_PHONE,
		Account:    "yanue",
		Phone:      "18503002165",
		VerifyCode: "589780",
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

func getUserInfo() {
	//fmt.Println("ctx", ctx)
	// Make request
	rsp, err := cl.GetUserInfo(context.Background(), &proto.PInt32{Val: 1})

	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("rsp:", rsp)
}

func getAccountInfo() {
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
