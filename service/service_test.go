/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountService.go
 Time    : 2018/9/11 15:16
 Author  : yanue
 Desc    : account微服务业务处理

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
	"github.com/yanue/go-esport-common/proto"
	"testing"
)

func TestAccountRpc_Login(t *testing.T) {
	var auth *AuthLogin = new(AuthLogin)
	auth.Auth = TUserAuth{
		AuthSite:    "qq",
		AuthOpenid:  "222",
		AuthUnionID: "222",
		AuthToken:   "",
		AuthExpire:  0,
	}
	auth.User = TUser{
		Name:   "yanue",
		Avatar: "aaa",
		Gender: 1,
	}

	u, err := rpc.account.setBindInfo(auth)
	fmt.Println("u, err", u, err)
}

func TestAccountService_Login(t *testing.T) {
	var device = &proto.PDevice{
		Imei: "19AAB430-9CB8-4325-ACC5-D7D386B68960",
		// 操作系统类型
		Os: proto.Os_IOS,
		// 操作系统版本
		OsVersion: "12.0.2",
		// 设备型号，如iPhone 6s
		Model: "iPhone 6s",
	}

	in := &proto.PLoginData{
		LoginType:  proto.ELoginType_PHONE,
		Phone:      "18503002165",
		VerifyCode: "929546",
		Device:     device,
	}

	// 处理登陆逻辑
	user, err := rpc.account.Login(in, "1.1.1.1")
	fmt.Println("user, err", user, err)
}

func TestAccountRpc_BindPhone(t *testing.T) {

	user, err := rpc.account.GetUserInfo(9)
	if err != nil {
		//return errcode.GetError(errcode.ErrAccountGetUserInfo, err.Error())
	}
	var auth *AuthLogin = new(AuthLogin)
	auth.Auth = TUserAuth{
		AuthSite:    "qq",
		AuthOpenid:  "222",
		AuthUnionID: "222",
		AuthToken:   "",
		AuthExpire:  0,
	}
	auth.User = TUser{
		Name:   "yanue",
		Avatar: "aaa",
		Gender: 1,
	}

	in := &proto.PBindData{
		BindType:      proto.ELoginType_QQ,
		QqOpenid:      "779E718A6B34B38807B83A7E2E649920",
		QqAccessToken: "CEF076B476582DD9C07F70DFE39DE802",
	}

	err = rpc.account.BindPhone(in, user)

	fmt.Println("err", err)
}

func TestAccountService_BindQQ(t *testing.T) {
	uid := 9
	var auth *AuthLogin = new(AuthLogin)
	auth.Auth = TUserAuth{
		AuthSite:    "wx",
		AuthOpenid:  "222",
		AuthUnionID: "222",
		AuthToken:   "",
		AuthExpire:  0,
	}
	auth.User = TUser{
		Name:   "yanue",
		Avatar: "aaa",
		Gender: 1,
	}
	err := rpc.account.saveAuthBind(uid, auth)
	fmt.Println("err", err)
}
