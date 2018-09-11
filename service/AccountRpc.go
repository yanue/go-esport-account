/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountRpc.go
 Time    : 2018/9/11 15:20
 Author  : yanue
 Desc    : account微服务rpc接口

------------------------------- go ---------------------------------*/

package service

import (
	"go-esport-account/proto"
	"golang.org/x/net/context"
)

/* 微服务接口 */
type AccountRpc struct {
	proto.AccountHandler
	account *AccountService
}

/**
注册
 */
func (this *AccountRpc) Reg(ctx context.Context, in *proto.PSingleString, out *proto.User) error {
	out.Name = "Reg"
	this.account.Reg()
	return nil
}

/**
登陆
 */
func (this *AccountRpc) Login(ctx context.Context, in *proto.PSingleString, out *proto.User) error {
	out.Name = "Login"
	return nil
}

/**
获取用户信息
 */
func (this *AccountRpc) GetUserInfo(ctx context.Context, in *proto.PInt32, out *proto.User) error {
	out.Name = "GetUserInfo"
	return nil
}

/**
获取验证码
 */
func (this *AccountRpc) GetVerifyCode(ctx context.Context, in *proto.PSingleString, out *proto.PBool) error {
	out.B = true
	return nil
}
