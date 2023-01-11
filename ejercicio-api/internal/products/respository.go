package products

import (
	"errors"
	"go-web-api/internal/domain"
)

var (
	ErrProdutNotFound = errors.New("Product doesn't exists")
)

type Repository interface {
	// read
	GetAllProducts() (products []domain.Product, err error)
	GetProductById(id int) (domain.Product, error)
	GetProductsMoreExpensiveThan(price float64) []domain.Product
	ExistCodeValue(codeValue string) (response bool)

	// write
	CreateProduct(domain.Product) (int, error)
	Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	Delete(id int) error
}

type repository struct {
	db     *[]domain.Product
	lastID int
}

func NewRepository(db *[]domain.Product, lastId int) Repository {
	return &repository{db: db, lastID: lastId}
}

// read functions
func (r *repository) GetAllProducts() ([]domain.Product, error) {
	return *r.db, nil
}

func (r *repository) GetProductById(id int) (domain.Product, error) {
	for _, p := range *r.db {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Product{}, ErrProdutNotFound
}

func (r *repository) GetProductsMoreExpensiveThan(price float64) []domain.Product {
	products := make([]domain.Product, 0)
	for _, p := range *r.db {
		if p.Price > price {
			products = append(products, p)
		}
	}
	return products
}

func (r *repository) ExistCodeValue(codeValue string) (response bool) {
	for _, product := range *r.db {
		if product.Code_value == codeValue {
			response = !response
			return
		}
	}
	return
}

// write functions
func (r *repository) CreateProduct(newProduct domain.Product) (int, error) {
	r.lastID++
	newProduct.ID = r.lastID
	*r.db = append(*r.db, newProduct)
	return r.lastID, nil
}

func (r *repository) Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error) {
	var updated bool
	products := *r.db
	var updatedProduct = domain.Product{
		Name:         name,
		Quantity:     quantity,
		Code_value:   code_value,
		Is_published: is_published,
		Expiration:   expiration,
		Price:        price,
	}

	for i, p := range products {
		if p.ID == id {
			updatedProduct.ID = id
			products[i] = updatedProduct
			updated = !updated
		}
	}

	if !updated {
		return domain.Product{}, ErrProdutNotFound
	}
	return updatedProduct, nil
}

func (r *repository) Delete(id int) (err error) {
	var deleted bool
	products := *r.db
	for index := range products {
		if products[index].ID != id {
			continue
		}
		products = append(products[:index], products[index+1:]...)
		deleted = !deleted
		break
	}
	if !deleted {
		err = ErrProdutNotFound
		return
	}
	return
}
