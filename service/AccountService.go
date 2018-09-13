/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountService.go
 Time    : 2018/9/11 15:16
 Author  : yanue
 Desc    : account微服务业务处理

------------------------------- go ---------------------------------*/

package service

import (
	"github.com/yanue/go-esport-common"
	"github.com/yanue/go-esport-common/proto"
)

type AccountService struct {
	orm   *AccountOrm
	cache *AccountCache
}

func (this *AccountService) Reg() {

}

func (this AccountService) Login(in *proto.PLoginData) (user *TUser, err error) {
	common.Logs.Info("")
	user = new(TUser)
	user.Phone = in.Phone
	user.Name = "yanue"
	user.Id = 1
	return user, nil
}
