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
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data interface{} `json:"data"`
}

func loadProducts(path string) ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func writeProducts(path string, list []domain.Product) error {
	bytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return err
}

func createServerForProductHandlerTest() *gin.Engine {
	godotenv.Load("../.env")
	_ = os.Setenv("TOKEN", "my-secret-token")
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

func createRequestTest(method string, url string, token string, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("token", token)

	return request, httptest.NewRecorder()
}

// En TestGetProducts_Ok se espera obtener todos los productos guardados
func TestGetProducts_OK(t *testing.T) {
	// Arrange
	r := createServerForProductHandlerTest()

	request, response := createRequestTest(http.MethodGet, "/products", "my-secret-token", "")

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

	request, response := createRequestTest(http.MethodGet, "/products/1", "my-secret-token", "")
	actual := map[string]domain.Product{}

	// Act
	r.ServeHTTP(response, request)

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, productExpected.Data, actual["data"])
}

// En TestCreate_Ok se espera obtener el producto creado como response
func TestCreate_Ok(t *testing.T) {
	// Arrange
	var expectd = response{Data: domain.Product{
		ID:           501,
		Name:         "Oil - Margarine",
		Quantity:     439,
		Code_value:   "TEST45050",
		Is_published: true,
		Expiration:   "15/12/2023",
		Price:        50.50,
	}}

	product, _ := json.Marshal(expectd.Data)

	r := createServerForProductHandlerTest()
	request, response := createRequestTest(http.MethodPost, "/products", "my-secret-token", string(product))

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		t.Fatal(err)
	}

	// Act
	r.ServeHTTP(response, request)
	actualProduct := map[string]domain.Product{}

	if err := writeProducts("./products_copy.json", p); err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusCreated, response.Code)
	err = json.Unmarshal(response.Body.Bytes(), &actualProduct)
	assert.Nil(t, err)
	assert.Equal(t, expectd.Data, actualProduct["data"])
}

// En TestDeleteOne_Ok se espera eliminar un producto del store
func TestDeleteOne_Ok(t *testing.T) {
	// Arrange
	r := createServerForProductHandlerTest()
	request, response := createRequestTest(http.MethodDelete, "/products/1", "my-secret-token", "")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		t.Fatal(err)
	}

	// Act
	r.ServeHTTP(response, request)

	err = writeProducts("./products_copy.json", p)
	if err != nil {
		panic(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
}

// En TestFail_ErrorBadRequest se obtiene un bad request y un mensaje informando del error
func TestFail_ErrorBadRequest(t *testing.T) {
	// Arrange
	test := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}
	r := createServerForProductHandlerTest()

	for _, tst := range test {
		request, response := createRequestTest(tst, "/products/aakdfn", "my-secret-token", "")
		r.ServeHTTP(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	}
}

// En TestFail_ErrorNotFound se obtiene un not found
func TestFail_ErrorNotFound(t *testing.T) {
	// Arrange
	test := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}
	r := createServerForProductHandlerTest()

	for _, tst := range test {
		request, response := createRequestTest(tst, "/products/505", "my-secret-token", "")
		r.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	}
}

// En TestFail_UnauthorizedError se obtiene un error 401
func TestFail_UnauthorizedError(t *testing.T) {
	// Arrange
	test := []string{http.MethodPut, http.MethodPatch, http.MethodDelete}
	r := createServerForProductHandlerTest()

	for _, tst := range test {
		request, response := createRequestTest(tst, "/products/334", "token1234", "")
		r.ServeHTTP(response, request)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	}
}
