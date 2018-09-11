package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	"go-esport-account/proto"
	"log"
)

var (
	cl proto.AccountClient
)

func InitApiService() {
	fmt.Println("InitApiService start")

	// Create service
	service := web.NewService(
		web.Name("go.esport.account.api"),
	)
	service.Init()

	// setup Greeter Server Client
	cl = proto.NewAccountClient("go.esport.account.srv", client.DefaultClient)

	// Create RESTful handler (using Gin)
	api := new(AccountApi)
	router := gin.Default()
	router.GET("/test", api.Anything)
	router.GET("/greeter/:name", api.Reg)

	// Register Handler
	service.Handle("/", router)

	go func() {
		// Run server
		if err := service.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("InitApiService end")
}
