package main

import (
	"encoding/json"
	"fmt"
	"go-web-api/handlers"
	"go-web-api/products"
	"os"

	"github.com/gin-gonic/gin"
)

func readProductJson() {
	data, err := os.ReadFile("products.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(data, &products.Products); err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	readProductJson()

	router := gin.Default()
	productsRouter := router.Group("/products")

	router.GET("/ping", handlers.Ping)
	productsRouter.GET("/", handlers.GetAllProducts)
	productsRouter.GET("/:id", handlers.GetProductById)
	productsRouter.GET("/search", handlers.GetProductsMoreExpensiveThan)

	router.Run(":8080")
}
