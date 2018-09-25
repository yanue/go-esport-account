/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountService.go
 Time    : 2018/9/11 15:16
 Author  : yanue
 Desc    : account微服务业务处理

------------------------------- go ---------------------------------*/

package service

import (
	"github.com/yanue/go-esport-common"
)

type AccountService struct {
	orm   *AccountOrm
	cache *AccountCache
}

func (this *AccountService) GetUserInfo(uid int) (user *TUser, err error) {
	// 从缓存中读取
	user, err = this.cache.GetUserInfo(uid)
	if err == nil {
		return user, nil
	}

	// 数据未找到,从mysql读取
	user, err = this.orm.GetUserInfo(uid)
	if err != nil {
		common.Logs.Info("user, err", user, err)
		return
	}

	this.cache.SetUserInfo(uid)

	return user, nil
}
