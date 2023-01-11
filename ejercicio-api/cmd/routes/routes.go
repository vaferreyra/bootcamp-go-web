package routes

import (
	"go-web-api/cmd/handlers"
	"go-web-api/internal/domain"
	"go-web-api/internal/product"

	"github.com/gin-gonic/gin"
)

type Router struct {
	db *[]domain.Product
	en *gin.Engine
}

func NewRouter(db *[]domain.Product, en *gin.Engine) *Router {
	return &Router{db: db, en: en}
}

func (r *Router) SetRoutes() {
	r.SetProduct()
}

func (r *Router) SetProduct() {
	rp := product.NewRepository(r.db, 500)
	sv := product.NewService(rp)
	h := handlers.NewProduct(sv)

	productRoute := r.en.Group("/products")

	// Read
	productRoute.GET("", h.GetAll())
	productRoute.GET("/:id", h.GetById())
	productRoute.GET("/search", h.GetMoreExpensiveThan())

	// Write
	productRoute.POST("", h.Create())
	productRoute.PUT("/:id", h.Update())

}
