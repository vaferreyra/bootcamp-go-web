package domain

type Product struct {
	ID           int     `json:"id" validator:"required"`
	Name         string  `json:"name" validator:"required"`
	Quantity     int     `json:"quantity" validator:"required"`
	Code_value   string  `json:"code_value" validator:"required"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration" validator:"required"`
	Price        float64 `json:"price" validator:"required"`
}
