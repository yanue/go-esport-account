package api

import (
	"github.com/gin-gonic/gin"
	"go-esport-account/proto"
	"golang.org/x/net/context"
	"log"
)

type AccountApi struct {
}

func (s *AccountApi) Anything(c *gin.Context) {
	log.Print("Received Say.Anything API request")
	c.JSON(200, map[string]string{
		"message": "Hi, this is the Greeter API",
	})
}

func (this *AccountApi) Reg(c *gin.Context) {
	log.Print("Received Say.Hello API request")

	name := c.Param("name")

	response, err := cl.Reg(context.TODO(), &proto.PSingleString{
		Str: name,
	})

	if err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, response)
}
