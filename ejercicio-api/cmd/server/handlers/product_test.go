package handlers

import (
	"bytes"
	"encoding/json"
	"go-web-api/internal/domain"
	"go-web-api/internal/products"
	"go-web-api/pkg/store"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data interface{} `json:"data"`
}

func createServerForProductHandlerTest() *gin.Engine {
	_ = os.Setenv("token", "my-secret-token")
	db := store.NewJSONStore("./products_copy.json")
	repo := products.NewRepository(db)
	service := products.NewService(repo)

	productService := NewProduct(service)
	router := gin.Default()

	productRoute := router.Group("/products")

	{
		// Read
		productRoute.GET("", productService.GetAll())
		productRoute.GET("/:id", productService.GetById())
		productRoute.GET("/search", productService.GetMoreExpensiveThan())

		// Write
		productRoute.POST("", productService.Create())
		productRoute.PUT("/:id", productService.Update())
		productRoute.PATCH("/:id", productService.PartialUpdate())
		productRoute.DELETE("/:id", productService.Delete())
	}

	return router
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("token", "my-secret-token")

	return request, httptest.NewRecorder()
}

// En TestGetProducts_Ok se espera obtener todos los productos guardados
func TestGetProducts_OK(t *testing.T) {
	// Arrange
	r := createServerForProductHandlerTest()

	request, response := createRequestTest(http.MethodGet, "/products", "")

	// Act
	r.ServeHTTP(response, request)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, len(body) > 0)
}

// En TestGetOne_Ok se espera obtener el producto con el id dado
func TestGetOne_Ok(t *testing.T) {
	// Arrange
	r := createServerForProductHandlerTest()
	productExpected := response{Data: domain.Product{
		ID:           1,
		Name:         "Oil - Margarine",
		Quantity:     439,
		Code_value:   "S82254D",
		Is_published: true,
		Expiration:   "15/12/2021",
		Price:        71.42,
	}}

	request, response := createRequestTest(http.MethodGet, "/products/1", "")
	actual := map[string]domain.Product{}

	// Act
	r.ServeHTTP(response, request)

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, productExpected.Data, actual["data"])
}
