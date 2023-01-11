package product

import (
	"errors"
	"go-web-api/internal/domain"
	"strconv"
	"strings"
)

var (
	ErrProductCodeAlreadyExist = errors.New("Code value already exists")
	ErrInvalidDate             = errors.New("Invalid date of expiration")
	ErrExpirationLength        = errors.New("Expiration date must have XX/XX/XXXX format")
	ErrExpirationNotNumber     = errors.New("Expiration date must be numbers")
)

type Service interface {
	// read
	GetAllProducts() (products []domain.Product, err error)
	GetProductById(id int) (domain.Product, error)
	GetProductsMoreExpensiveThan(price float64) []domain.Product

	// write
	CreateProduct(name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
}

type service struct {
	rp Repository
}

func NewService(rp Repository) Service {
	return &service{rp: rp}
}

func (service *service) GetAllProducts() (products []domain.Product, err error) {
	return service.rp.GetAllProducts()
}

func (service *service) GetProductById(id int) (domain.Product, error) {
	return service.rp.GetProductById(id)
}

func (service *service) GetProductsMoreExpensiveThan(price float64) []domain.Product {
	return service.rp.GetProductsMoreExpensiveThan(price)
}

func (service *service) CreateProduct(name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error) {
	if service.rp.ExistCodeValue(code_value) {
		return domain.Product{}, ErrProductCodeAlreadyExist
	}

	_, err := IsValidExpiration(expiration)
	if err != nil {
		return domain.Product{}, err
	}

	newProduct := domain.Product{
		Name:         name,
		Quantity:     quantity,
		Code_value:   code_value,
		Is_published: is_published,
		Expiration:   expiration,
		Price:        price,
	}
	lastId, err := service.rp.CreateProduct(newProduct)
	if err != nil {
		return domain.Product{}, err
	}
	newProduct.ID = lastId
	return newProduct, nil
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
