package products

type Product struct {
	ID           int
	Name         string
	Quantity     int
	Code_value   int
	Is_published bool
	Expiration   string
	Price        float64
}

type ProductCatalog struct {
	Products []Product
}
