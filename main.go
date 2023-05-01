package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubenscholle/venus/authbundle"
	core "github.com/rubenscholle/venus/corebundle"
)

func main() {
	log.SetFlags(log.LstdFlags)

	err := core.InitConfig()
	if err != nil {
		log.Println(err)
	}

	core.OrmDb = core.InitDb()

	log.Println("starting server...")

	router := gin.Default()

	authbundle.InitBundle(router.Group("/auth"), core.OrmDb)

	router.GET("/hello-world", helloWorldHandler)
	router.Run(":7901")
}

func helloWorldHandler(c *gin.Context) {
	c.String(http.StatusOK, "hello world!\nthis is a RESTful API by Ruben Scholle")
}
