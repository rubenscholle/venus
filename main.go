package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello-world", helloWorldHandler)
	router.Run(":7901")
}

func helloWorldHandler(c *gin.Context) {
	c.String(http.StatusOK, "hello world!")
}
