package main

import (
	"encoding/json"
	"go-web-api/cmd/routes"
	"go-web-api/internal/domain"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadProducts(path string, list *[]domain.Product) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(file), &list)
	if err != nil {
		panic(err)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var db []domain.Product

	loadProducts("products.json", &db)

	en := gin.Default()
	router := routes.NewRouter(&db, en)
	router.SetRoutes()

	if err := en.Run(); err != nil {
		log.Fatal(err)
	}
}
