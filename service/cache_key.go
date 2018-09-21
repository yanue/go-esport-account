/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : redis key

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
)

type CRedisKey struct{}

/**
@note 用户信息
 */
func (this *CRedisKey) HUserInfo(uid int) (string) {
	return RedisPrefix + fmt.Sprintf("user:info:%d", uid)
}

/**
@note 用户信息
 */
func (this *CRedisKey) SUserToken(uid int) (string) {
	return RedisPrefix + fmt.Sprintf("user:token:%d", uid)
}
