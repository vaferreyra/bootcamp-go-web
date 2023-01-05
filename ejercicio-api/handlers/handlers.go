package handlers

import (
	"go-web-api/products"
	"go-web-api/productsService"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func GetAllProducts(c *gin.Context) {
	products := productsService.ProductsCatalog
	c.JSON(http.StatusOK, Response{Message: "Succeed to get all products", Data: products})
}

func GetProductById(c *gin.Context) {
	paramId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "Bad request error"})
		return
	}

	var productToReturn products.Product

	for _, product := range productsService.ProductsCatalog.Products {
		if product.ID == paramId {
			productToReturn = product
			break
		}
	}

	if productToReturn.ID != 0 {
		c.JSON(http.StatusOK, Response{Message: "Succeed to get product by id", Data: productToReturn})
		return
	} else {
		c.JSON(http.StatusNotFound, Response{Message: "Error to get product by id"})
		return
	}
}

func GetProductsMoreExpensiveThan(c *gin.Context) {
	price, err := strconv.ParseFloat(c.Query("priceGt"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Message: "Parameter invalid"})
		return
	}
	productsToReturn := make([]products.Product, 0)
	for _, product := range productsService.ProductsCatalog.Products {
		if price != 0 && product.Price >= price {
			productsToReturn = append(productsToReturn, product)
		}
	}
	c.JSON(http.StatusOK, Response{Message: "Succeed to get products", Data: productsToReturn})
}
