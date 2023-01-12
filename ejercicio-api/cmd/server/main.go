package main

import (
	"go-web-api/cmd/routes"
	"go-web-api/pkg/store"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	productStorage := store.NewJSONStore("products.json")

	en := gin.Default()
	router := routes.NewRouter(productStorage, en)
	router.SetRoutes()

	if err := en.Run(); err != nil {
		log.Fatal(err)
	}
}
