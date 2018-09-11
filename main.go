/* -------------- Copyright (c) Shenzhen BB Team. ------------------

 File    : main.go
 Time    : 2018/9/11 15:13
 Author  : yanue
 Desc    : 服务初始化

------------------------------- go --------------------------------- */

package main

import (
	"go-esport-account/service"
)

func main() {
	// 初始化
	service.InitAccountService()

	// 卸载
	service.UnInitAccountService()
}
