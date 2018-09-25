/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountRpc.go
 Time    : 2018/9/11 15:20
 Author  : yanue
 Desc    : account微服务rpc接口

------------------------------- go ---------------------------------*/

package service

import (
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

func (this *AccountRpc) verifyToken(ctx context.Context) (token *proto.PJwtToken, err error) {
	// 获取token
	md, ok := metadata.FromContext(ctx)
	tokenStr, ok := md["authorization"]
	if !ok {
		errcode.GetError(errcode.ErrInvalidParam, "authorization")
		return
	}

	// 验证token
	jwt, payload1, err := util.JwtToken.Verify(tokenStr)
	if err != nil {
		common.Logs.Info("Verify JwtToken err:", err.Error())
		err = errcode.GetError(errcode.ErrAccountTokenVerify)
		return
	}

	// 验证redis内token信息
	payload2, err := this.account.cache.GetUserToken(int(jwt.Payload.Uid))
	if err != nil {
		common.Logs.Info("Verify JwtToken GetUserToken:", err.Error())
		err = errcode.GetError(errcode.ErrAccountTokenGet)
		return
	}

	if payload1 != payload2 {
		err = errcode.GetError(errcode.ErrAccountTokenNotEqual)
		return
	}

	// 返回token
	return jwt.PJwtToken, nil
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
	// 取地址
	*out = *user

	return nil
}

/**
获取用户信息
 */
func (this *AccountRpc) GetAccountInfo(ctx context.Context, in *proto.PNoParam, out *proto.PUser) error {
	token, err := this.verifyToken(ctx)
	if err != nil {
		return err
	}

	user, err := this.account.GetUserInfo(int(token.Payload.Uid))
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	// 取地址
	*out = proto.PUser{
		Id:             int32(user.Id),
		Name:           user.Name,
		Phone:          user.Phone,
		Email:          user.Email,
		Gender:         int32(user.Gender),
		SchoolId:       int32(user.SchoolId),
		ClassId:        int32(user.ClassId),
		AreaId:         int32(user.AreaId),
		IdentityStatus: int32(user.IdentityStatus),
		Created:        int32(user.Created),
	}

	return nil
}

/**
获取用户信息
 */
func (this *AccountRpc) GetUserInfo(ctx context.Context, in *proto.PInt32, out *proto.PUser) error {
	user, err := this.account.GetUserInfo(int(in.Val))
	if err != nil {
		return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}

	// 取地址
	*out = proto.PUser{
		Id:             int32(user.Id),
		Name:           user.Name,
		Phone:          user.Phone,
		Email:          user.Email,
		Gender:         int32(user.Gender),
		SchoolId:       int32(user.SchoolId),
		ClassId:        int32(user.ClassId),
		AreaId:         int32(user.AreaId),
		IdentityStatus: int32(user.IdentityStatus),
		Created:        int32(user.Created),
	}

	return nil
}

/**
获取验证码
 */
func (this *AccountRpc) GetVerifyCode(ctx context.Context, in *proto.PString, out *proto.PBool) error {
	out.Val = true
	return nil
}
