package main

import (
	"go-web-api/cmd/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	productsRouter := router.Group("/products")

	//------- GET -------
	router.GET("/ping", handlers.Ping)
	productsRouter.GET("", handlers.GetAllProducts)
	productsRouter.GET("/:id", handlers.GetProductById)
	productsRouter.GET("/search", handlers.GetProductsMoreExpensiveThan)

	//------- POST -------
	productsRouter.POST("", handlers.CreateProduct)

	router.Run(":8080")
}
