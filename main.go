package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rubenscholle/venus/authbundle"
	core "github.com/rubenscholle/venus/corebundle"
	"github.com/rubenscholle/venus/websocketbundle"
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
	websocketbundle.InitBundle(publicRoutes, protectedRoutes.Group("websocket"), core.OrmDb)

	protectedRoutes.GET("/version", versionHandler)
	router.Run(fmt.Sprintf(":%d", core.Config.Server.Port))
}

func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "v0.1.1c")
}
