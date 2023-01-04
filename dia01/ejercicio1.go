package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main_1() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
