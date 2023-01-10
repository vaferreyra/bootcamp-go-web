package handlers

import (
	"errors"
	"go-web-api/internal/product"
	"go-web-api/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidParameter = errors.New("Invalid parameter, try another")
)

type Product struct {
	sv product.Service
}

func NewProduct(sv product.Service) *Product {
	return &Product{sv: sv}
}

func (p *Product) GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := p.sv.GetAllProducts()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
		ctx.JSON(http.StatusOK, response.Ok("Success to get products", products))
	}
}

func (p *Product) GetProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paramId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidParameter))
			return
		}

		product, err := p.sv.GetProductById(paramId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, response.Err(err))
			return
		}

		ctx.JSON(http.StatusOK, response.Ok("Success to get product", product))
	}
}

func (p *Product) GetProductsMoreExpensiveThan() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidParameter))
			return
		}
		products := p.sv.GetProductsMoreExpensiveThan(price)
		ctx.JSON(http.StatusOK, response.Ok("Success to get products", products))
	}
}

func (p *Product) CreateProduct() gin.HandlerFunc {
	type request struct {
		Name         string  `json:"name" validate:"required"`
		Quantity     int     `json:"quantity" validate:"required"`
		Code_value   string  `json:"code_value" validate:"required"`
		Is_published bool    `json:"is_published"`
		Expiration   string  `json:"expiration" validate:"required"`
		Price        float64 `json:"price" validate:"required"`
	}

	return func(ctx *gin.Context) {
		var req request

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(err))
			return
		}

		product, err := p.sv.CreateProduct(req.Name, req.Quantity, req.Code_value, req.Is_published, req.Expiration, req.Price)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.Err(err))
			return
		}

		ctx.JSON(http.StatusCreated, response.Ok("Success to create product", product))
	}
}
