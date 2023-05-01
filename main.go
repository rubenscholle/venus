package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubenscholle/venus/authbundle"
	core "github.com/rubenscholle/venus/corebundle"
	"gorm.io/gorm"
)

func main() {
	log.SetFlags(log.LstdFlags)

	err := core.InitConfig()
	if err != nil {
		log.Println(err)
	}

	core.OrmDb = core.InitDb()

	// ToDo move to controllers/bundles
	AutoMigrate(core.OrmDb)

	log.Println("starting server...")

	router := gin.Default()
	router.GET("/hello-world", helloWorldHandler)

	authController := authbundle.AuthController{OrmDb: *core.OrmDb}
	router.POST("/login", authController.LoginHandler)

	router.Run(":7901")
}

func helloWorldHandler(c *gin.Context) {
	c.String(http.StatusOK, "hello world!\nthis is a RESTful API by Ruben Scholle")
}

func AutoMigrate(ormDb *gorm.DB) {
	ormDb.AutoMigrate(&authbundle.AuthUser{})
}
