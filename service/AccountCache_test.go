/*-------------- Copyright (c) Shenzhen BB Team. -------------------

 File    : AccountBase.go
 Time    : 2018/9/11 15:19
 Author  : yanue
 Desc    : account微服务-redis操作

------------------------------- go ---------------------------------*/

package service

import (
	"fmt"
	"testing"
)

func init() {
}

func TestAccountCache_GetUserInfo(t *testing.T) {
	res, err := cache.GetUserInfo(1)
	fmt.Println("res, err", res, err)
}
