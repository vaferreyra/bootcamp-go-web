package handlers

import (
	"encoding/json"
	"errors"
	"go-web-api/internal/products"
	product "go-web-api/internal/products"
	"go-web-api/pkg/web"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidParameter    = errors.New("Invalid parameter, try another")
	ErrInvalidId           = errors.New("Id invalid")
	ErrInvalidBody         = errors.New("Something's wrong with the body")
	ErrInvalidDate         = errors.New("Invalid date of expiration")
	ErrExpirationLength    = errors.New("Expiration date must have XX/XX/XXXX format")
	ErrExpirationNotNumber = errors.New("Expiration date must be numbers")
	ErrEmptyName           = errors.New("The product's name cannot be empty")
	ErrEmptyExpiration     = errors.New("The product's expiration cannot be empty")
	ErrEmptyCodeValue      = errors.New("The product's code value cannot be empty")
	ErrInvalidQuantity     = errors.New("The product's quantity must be > 0")
	ErrInvalidPrice        = errors.New("The product's price must be > 0")
	ErrUserUnauthorized    = errors.New("User unauthorized")
)

type NewProductRequest struct {
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

type Product struct {
	sv product.Service
}

func NewProduct(sv product.Service) *Product {
	return &Product{sv: sv}
}

func (p *Product) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := p.sv.GetAllProducts()
		if err != nil {
			web.Failure(ctx, http.StatusInternalServerError, nil)
			return
		}
		web.Success(ctx, http.StatusOK, products)
	}
}

func (p *Product) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paramId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}

		product, err := p.sv.GetProductById(paramId)
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, nil)
			return
		}

		web.Success(ctx, http.StatusOK, product)
	}
}

func (p *Product) GetMoreExpensiveThan() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}
		products := p.sv.GetProductsMoreExpensiveThan(price)
		web.Success(ctx, http.StatusOK, products)
	}
}

func (p *Product) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, nil)
			return
		}

		var req NewProductRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}

		_, err := IsValidProduct(req)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}

		product, err := p.sv.CreateProduct(req.Name, req.Quantity, req.Code_value, req.Is_published, req.Expiration, req.Price)
		if err != nil {
			web.Failure(ctx, http.StatusInternalServerError, nil)
			return
		}

		web.Success(ctx, http.StatusCreated, product)
	}
}

func (p *Product) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Consigo token de usuario y verifico que sea valido
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, nil)
			return
		}

		// Busco si el producto existe en DB
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}
		var request NewProductRequest

		if err := ctx.ShouldBindJSON(&request); err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}

		_, err = IsValidProduct(request)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}

		result, err := p.sv.Update(int(id), request.Name, request.Quantity, request.Code_value, request.Is_published, request.Expiration, request.Price)
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, nil)
			return
		}

		web.Success(ctx, http.StatusCreated, result)
	}
}

func (p *Product) PartialUpdate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtengo token de usuario y verifico que sea valido
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, nil)
			return
		}

		// Obtengo el id pasado por parámetro
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}

		// Me traigo el producto desde la base
		product, err := p.sv.GetProductById(int(id))
		if err != nil {
			switch err {
			case products.ErrProdutNotFound:
				web.Failure(ctx, http.StatusNotFound, nil)
			default:
				ctx.JSON(http.StatusInternalServerError, nil)
			}
			return
		}
		if err = json.NewDecoder(ctx.Request.Body).Decode(&product); err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}
		productUpdated, err := p.sv.Update(product.ID, product.Name, product.Quantity, product.Code_value, product.Is_published, product.Expiration, product.Price)
		if err != nil {
			switch err {
			case products.ErrProdutNotFound:
				web.Failure(ctx, http.StatusNotFound, nil)
			default:
				ctx.JSON(http.StatusInternalServerError, nil)
			}
			return
		}

		web.Success(ctx, http.StatusOK, productUpdated)
	}
}

func (p *Product) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Obtengo token de usuario y verifico que sea valido
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, nil)
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, nil)
			return
		}

		if err := p.sv.Delete(int(id)); err != nil {
			web.Failure(ctx, http.StatusNotFound, nil)
			return
		}

		web.Success(ctx, http.StatusOK, id)
	}
}

func IsValidProduct(product NewProductRequest) (result bool, err error) {
	if product.Name == "" {
		err = ErrEmptyName
		return
	}

	if product.Quantity <= 0 {
		err = ErrInvalidQuantity
		return
	}

	if product.Code_value == "" {
		err = ErrEmptyCodeValue
		return
	}

	if product.Price <= 0 {
		err = ErrInvalidPrice
		return
	}

	if product.Expiration == "" {
		err = ErrEmptyExpiration
		return
	}

	_, er := IsValidExpiration(product.Expiration)
	if er != nil {
		err = er
		return
	}

	result = !result
	return
}

func IsValidExpiration(date string) (result bool, err error) {
	dateFormatted := strings.Split(date, "/")
	listOfInt := []int{}
	if len(dateFormatted) != 3 {
		err = ErrExpirationLength
		return
	}

	for _, value := range dateFormatted {
		v, er := strconv.Atoi(value)
		if er != nil {
			err = ErrExpirationNotNumber
			return
		}
		listOfInt = append(listOfInt, v)
	}

	day := listOfInt[0]
	month := listOfInt[1]
	year := listOfInt[2]

	condition := day > 0 && day <= 31 && month > 0 && month <= 12 && year > 2022
	if !condition {
		err = ErrInvalidDate
		return
	}
	result = !result
	return
}
