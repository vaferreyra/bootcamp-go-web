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
		var contextName FullName

		if err := ctx.BindJSON(&contextName); err != nil {
			panic(err)
		}

		ctx.String(http.StatusOK, "Hola %s %s", contextName.Nombre, contextName.Apellido)
	})

	router.Run(":8080")
}
