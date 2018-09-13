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
	// todo
	// 读取配置信息
	RdsUser := "root"
	RdsAuth := "blemobi3721"
	RdsAddr := "192.168.5.202"
	RdsDatabase := "esport-account"
	RedisDsn := "192.168.5.201:6379"
	RedisPass := ""

	// 初始化
	service.InitAccountService(RdsUser, RdsAuth, RdsAddr, RdsDatabase, RedisDsn, RedisPass)
	// 卸载
	service.UnInitAccountService()
}
