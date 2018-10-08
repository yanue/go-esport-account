/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountService.go
 Time    : 2018/9/11 15:16
 Author  : yanue
 Desc    : account微服务业务处理

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
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
