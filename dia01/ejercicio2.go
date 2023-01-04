package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FullName struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

func main() {
	router := gin.Default()

	router.POST("/saludo", func(ctx *gin.Context) {
		var r FullName

		if err := ctx.BindJSON(&r); err != nil {
			panic(err)
		}

		ctx.String(http.StatusOK, "Hola %s %s", r.Nombre, r.Apellido)

		// ctx.String(http.StatusOK, "Hola %s %s", nombre, apellido)
	})

	router.Run(":8080")
}
