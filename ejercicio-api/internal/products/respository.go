package products

import (
	"errors"
	"go-web-api/internal/domain"
	"go-web-api/pkg/store"
)

var (
	ErrProdutNotFound = errors.New("Product doesn't exists")
)

type Repository interface {
	// read
	GetAllProducts() (products []domain.Product, err error)
	GetProductById(id int) (domain.Product, error)
	GetProductsMoreExpensiveThan(price float64) []domain.Product

	// write
	CreateProduct(domain.Product) (int, error)
	Update(id int, product domain.Product) (domain.Product, error)
	Delete(id int) error

	// private
	existCodeValue(codeValue string) (response bool)
}

type repository struct {
	db store.Store
}

// NewRepository crea un nuevo Repository que interactuará con la db dada
func NewRepository(db store.Store) Repository {
	return &repository{db}
}

// read functions
func (r *repository) GetAllProducts() (ps []domain.Product, er error) {
	products, err := r.db.GetAll()
	if err != nil {
		er = err
		return
	}
	ps = products
	return
}

func (r *repository) GetProductById(id int) (p domain.Product, er error) {
	product, err := r.db.GetOne(id)
	if err != nil {
		er = err
		return
	}
	p = product
	return
}

func (r *repository) GetProductsMoreExpensiveThan(price float64) []domain.Product {
	var products []domain.Product
	list, err := r.db.GetAll()
	if err != nil {
		return products
	}

	for _, p := range list {
		if p.Price > price {
			products = append(products, p)
		}
	}
	return products
}

func (r *repository) existCodeValue(codeValue string) (response bool) {
	list, err := r.db.GetAll()
	if err != nil {
		return
	}

	for _, product := range list {
		if product.Code_value == codeValue {
			response = !response
			return
		}
	}
	return
}

// write functions
func (r *repository) CreateProduct(newProduct domain.Product) (id int, err error) {
	if r.existCodeValue(newProduct.Code_value) {
		err = ErrProductCodeAlreadyExist
		return
	}

	idCreated, er := r.db.AddOne(newProduct)
	if er != nil {
		err = er
		return
	}

	id = idCreated
	return
}

func (r *repository) Update(id int, p domain.Product) (product domain.Product, er error) {
	if r.existCodeValue(product.Code_value) {
		er = ErrProductCodeAlreadyExist
		return
	}
	p.ID = id
	if err := r.db.UpdateOne(p); err != nil {
		er = err
		return
	}
	product = p
	return
}

func (r *repository) Delete(id int) (er error) {
	if err := r.db.DeleteOne(id); err != nil {
		er = err
		return
	}
	return
}
