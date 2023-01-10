package handlers

import (
	"errors"
	"go-web-api/pkg/response"
	"go-web-api/services"
	"go-web-api/services/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong")
}

func GetAllProducts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.Ok("Succeed to get all products", services.ProductsCatalog.Products))
}

func GetProductById(ctx *gin.Context) {
	paramId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return
	}

	var productToReturn models.Product

	for _, product := range services.ProductsCatalog.Products {
		if product.ID == paramId {
			productToReturn = product
			break
		}
	}

	if productToReturn.ID != 0 {
		ctx.JSON(http.StatusOK, response.Ok("Succeed to get product by id", productToReturn))
		return
	} else {
		ctx.JSON(http.StatusNotFound, response.Err(errors.New("Error to get product by id")))
		return
	}
}

func GetProductsMoreExpensiveThan(ctx *gin.Context) {
	price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return
	}
	productsToReturn := make([]models.Product, 0)
	for _, product := range services.ProductsCatalog.Products {
		if price != 0 && product.Price >= price {
			productsToReturn = append(productsToReturn, product)
		}
	}
	ctx.JSON(http.StatusOK, response.Ok("Succeed to get products", productsToReturn))
}

func CreateProduct(ctx *gin.Context) {
	var request models.RequestProduct

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return
	}

	validate := validator.New()

	if err := validate.Struct(request); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, response.Err(err))
		return
	}

	newProduct, err := services.Create(request.Name, request.Quantity, request.Code_value, request.Is_published, request.Expiration, request.Price)
	switch err {
	case services.ErrProductCodeAlreadyExist:
		ctx.JSON(http.StatusConflict, response.Err(err))
		return
	case services.ErrInternalError:
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}
	ctx.JSON(http.StatusCreated, response.Ok("Product created successfuly", newProduct))
}
