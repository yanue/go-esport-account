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

var device = &proto.PDevice{
	Imei: "19AAB430-9CB8-4325-ACC5-D7D386B68960",
	// 操作系统类型
	Os: proto.Os_IOS,
	// 操作系统版本
	OsVersion: "12.0.2",
	// 设备型号，如iPhone 6s
	Model: "iPhone 6s",
}

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

func loginByAccount() {
	pass, _ := util.Rsa.RsaEncryptPublic("111111")

	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		LoginType: proto.ELoginType_ACCOUNT,
		Account:   "yanue",
		Password:  pass,
		Device:    device,
	})

	if err != nil {
		fmt.Println("err:", err)
		return
	}
	token = rsp.Token
	fmt.Println("rsp:", common.OS_IOS, rsp)
}

func loginByPhone() {
	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		LoginType:  proto.ELoginType_PHONE,
		Phone:      "18503002165",
		VerifyCode: "929546",
		Device:     device,
	})

	if err != nil {
		fmt.Println("err:", err)
		return
	}
	token = rsp.Token
	fmt.Println("rsp:", common.OS_IOS, rsp)
}
func loginByQQ() {
	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		LoginType:     proto.ELoginType_QQ,
		QqOpenid:      "779E718A6B34B38807B83A7E2E649920",
		QqAccessToken: "CEF076B476582DD9C07F70DFE39DE802",
		Device:        device,
	})

	if err != nil {
		fmt.Println("err:", err)
		return
	}
	token = rsp.Token
	fmt.Println("rsp:", common.OS_IOS, rsp)
}

func loginByWechat() {
	// Make request
	rsp, err := cl.Login(context.Background(), &proto.PLoginData{
		LoginType: proto.ELoginType_WECHAT,
		WxCode:    "001PLaZK1k10U40ol9ZK1EanZK1PLaZb",
		Device:    device,
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
