package main

import (
	"github.com/micro/go-micro"
	"log"
)

func main() {
	service := micro.NewService(micro.Name("go-esport-account"))

	service.Init()

	err := service.Run()
	if err != nil {
		log.Fatalf("启动失败:%v", err.Error())
	}
}
