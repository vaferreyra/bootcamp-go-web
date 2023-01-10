package main

import (
	"encoding/json"
	"go-web-api/cmd/handlers"
	"go-web-api/services"
	"go-web-api/services/models"
	"os"

	"github.com/gin-gonic/gin"
)

func loadProducts(path string, list *[]models.Product) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(file), &list)
	if err != nil {
		panic(err)
	}
}

func main() {

	loadProducts("../products.json", &services.ProductsCatalog.Products)

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
