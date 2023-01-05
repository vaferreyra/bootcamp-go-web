package main

import (
	"go-web-api/handlers"
	"go-web-api/productsService"

	"github.com/gin-gonic/gin"
)

func main() {
	productsService.ReadProductJson()

	router := gin.Default()
	productsRouter := router.Group("/products")

	router.GET("/ping", handlers.Ping)
	productsRouter.GET("/", handlers.GetAllProducts)
	productsRouter.GET("/:id", handlers.GetProductById)
	productsRouter.GET("/search", handlers.GetProductsMoreExpensiveThan)

	router.Run(":8080")
}
