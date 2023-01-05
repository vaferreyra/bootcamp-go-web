package handlers

import (
	"go-web-api/products"
	"go-web-api/productsService"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func GetAllProducts(c *gin.Context) {
	products := productsService.ProductsCatalog
	c.JSON(http.StatusOK, gin.H{
		"message": "Succed to get all products",
		"data":    products.Products,
	})
}

func GetProductById(c *gin.Context) {
	paramId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"data":    nil,
		})
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
		c.JSON(http.StatusOK, gin.H{
			"message": "Succeed to find product by id",
			"data":    productToReturn,
		})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Error on find product by id",
			"data":    nil,
		})
		return
	}
}

func GetProductsMoreExpensiveThan(c *gin.Context) {
	price, err := strconv.ParseFloat(c.Query("priceGt"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Parameter invalid",
		})
		return
	}
	productsToReturn := make([]products.Product, 0)
	for _, product := range productsService.ProductsCatalog.Products {
		if price != 0 && product.Price >= price {
			productsToReturn = append(productsToReturn, product)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Succeed to filter products",
		"data":    productsToReturn,
	})
}
