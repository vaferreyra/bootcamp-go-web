package productsService

import (
	"encoding/json"
	"fmt"
	"go-web-api/products"
	"os"
)

var ProductsCatalog = products.ProductCatalog{}

func ReadProductJson() {
	data, err := os.ReadFile("./products.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(data, &ProductsCatalog.Products); err != nil {
		fmt.Println(err)
		return
	}
}
