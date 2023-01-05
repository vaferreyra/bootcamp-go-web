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

	router.GET("/ping", handlers.Ping)
	router.GET("/products", handlers.GetAllProducts)
	router.GET("/products/:id", handlers.GetProductById)
	router.GET("products/search", handlers.GetProductsMoreExpensiveThan)

	router.Run()
}
