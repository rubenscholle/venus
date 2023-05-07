package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
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

	// Add CORS support to all routes
	router.Use(
		cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch, http.MethodHead, http.MethodOptions},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
		}),
	)

	publicRoutes := router.Group("/auth")
	protectedRoutes := router.Group("/api/v1")
	protectedRoutes.Use(authbundle.Middleware())

	authbundle.InitBundle(publicRoutes, protectedRoutes.Group("auth"), core.OrmDb)

	protectedRoutes.GET("/hello-world", helloWorldHandler)
	router.Run(":7901")
}

func helloWorldHandler(c *gin.Context) {
	c.String(http.StatusOK, "hello world!\nthis is a RESTful API by Ruben Scholle")
}
