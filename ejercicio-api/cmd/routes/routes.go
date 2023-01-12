package routes

import (
	"go-web-api/cmd/server/handlers"
	product "go-web-api/internal/products"
	"go-web-api/pkg/store"

	"github.com/gin-gonic/gin"
)

type Router struct {
	db store.Store
	en *gin.Engine
}

func NewRouter(db store.Store, en *gin.Engine) *Router {
	return &Router{db: db, en: en}
}

func (r *Router) SetRoutes() {
	r.SetProduct()
}

func (r *Router) SetProduct() {
	rp := product.NewRepository(r.db)
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
	productRoute.PATCH("/:id", h.PartialUpdate())
	productRoute.DELETE("/:id", h.Delete())

}
