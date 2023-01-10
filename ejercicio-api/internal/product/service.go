package product

import (
	"errors"
	"go-web-api/internal/domain"
)

var (
	ErrProductCodeAlreadyExist = errors.New("Code value already exists")
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
