/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountRpc.go
 Time    : 2018/9/11 15:20
 Author  : yanue
 Desc    : account微服务rpc接口

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
	"github.com/micro/go-micro/metadata"
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/errcode"
	"github.com/yanue/go-esport-common/proto"
	"github.com/yanue/go-esport-common/util"
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
func (this *AccountRpc) Reg(ctx context.Context, in *proto.PString, out *proto.PUserAndToken) error {
	return nil
}

/**登陆*/
func (this *AccountRpc) Login(ctx context.Context, in *proto.PLoginData, out *proto.PUserAndToken) error {
	// 处理登陆逻辑
	user, err := this.account.Login(in)
	if err != nil {
		common.Logs.Debug("login err=", err.Error(), user)
		return err
	}

	// 不能直接赋值
	out.User = user.User
	out.Token = user.Token
	common.Logs.Info("user ", out)

	return nil
}

/**
获取用户信息
 */
func (this *AccountRpc) GetAccountInfo(ctx context.Context, in *proto.PInt32, out *proto.PUser) error {
	md, ok := metadata.FromContext(ctx)
	token, ok := md["authorization"]
	if !ok {
		return errcode.GetError(errcode.ErrAccountTokenVerify)
	}
	fmt.Println("token", token)
	//
	// 验证token
	_, err := util.JwtToken.Verify(token)
	if err != nil {
		common.Logs.Info("Verify err ", err.Error())
		return errcode.GetError(errcode.ErrAccountTokenVerify)
	}
	//
	////
	////fmt.Println("md, ok", md, ok)
	////fmt.Println("ctx", ctx)
	//user, err := this.account.GetUserInfo(int(in.Val))
	//if err != nil {
	//	return errcode.GetError(errcode.ErrAccountNotExist)
	//}
	//fmt.Println("user", user)
	//out.Name = user.Name
	//out.Email = user.Email
	//out.Phone = user.Phone

	return nil
}

/**
获取验证码
 */
func (this *AccountRpc) GetVerifyCode(ctx context.Context, in *proto.PString, out *proto.PBool) error {
	out.Val = true
	return nil
}
