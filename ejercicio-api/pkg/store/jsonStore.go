package store

import (
	"encoding/json"
	"errors"
	"go-web-api/internal/domain"
	"os"
)

var (
	ErrNotFound = errors.New("Product not found")
)

type Store interface {
	GetAll() ([]domain.Product, error)
	GetOne(id int) (domain.Product, error)
	AddOne(product domain.Product) (int, error)
	UpdateOne(product domain.Product) error
	DeleteOne(id int) error
	saveProducts(products []domain.Product) error
	loadProducts() ([]domain.Product, error)
}

type jsonStore struct {
	fileName string
}

// loadProducts carga los productso desde un archivo json
func (s *jsonStore) loadProducts() (products []domain.Product, er error) {
	file, err := os.ReadFile(s.fileName)
	if err != nil {
		er = err
		return
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		er = err
		return
	}
	return
}

// saveProducts guarda los productos en un archivo json
func (s *jsonStore) saveProducts(products []domain.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return os.WriteFile(s.fileName, bytes, 0644)
}

// NewJSONStore crea una nueva instancia de store de products
func NewJSONStore(file string) Store {
	return &jsonStore{fileName: file}
}

// GetAll devuelve todos los productos que se encuentran en el archivo
func (s *jsonStore) GetAll() ([]domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetOne devuelve un solo producto cuyo id coincida con el dado
func (s *jsonStore) GetOne(id int) (product domain.Product, er error) {
	ps, err := s.loadProducts()
	if err != nil {
		er = err
		return
	}
	for _, p := range ps {
		if p.ID == id {
			product = p
			return
		}
	}
	er = ErrNotFound
	return
}

// AddOne agrega un producto al store y lo guarda en el archivo
// devuelve el id del producto almacenado
func (s *jsonStore) AddOne(product domain.Product) (id int, er error) {
	ps, err := s.loadProducts()
	if err != nil {
		er = err
		return
	}
	product.ID = len(ps) + 1
	ps = append(ps, product)
	if err := s.saveProducts(ps); err != nil {
		er = err
		return
	}
	id = product.ID
	return
}

// UpdateOne cambia los datos del producto existente por los nuevos dados
func (s *jsonStore) UpdateOne(product domain.Product) error {
	ps, err := s.loadProducts()
	if err != nil {
		return err
	}

	for i, p := range ps {
		if p.ID == product.ID {
			ps[i] = product
			return s.saveProducts(ps)
		}
	}
	return ErrNotFound
}

// DeleteOne elimina un producto del store
func (s *jsonStore) DeleteOne(id int) error {
	ps, err := s.loadProducts()
	if err != nil {
		return err
	}

	for i, p := range ps {
		if p.ID == id {
			ps = append(ps[:i], ps[i+1:]...)
			return s.saveProducts(ps)
		}
	}

	return ErrNotFound
}
