package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/cityhunteur/data-service/internal/database"
)

func TestDataHandler_PostGetData(t *testing.T) {
	t.Parallel()

	testDB, _ := database.NewTestDatabase(t)

	h := &DataHandler{
		db: database.NewDataDB(testDB),
	}

	router := setupRouter(h)

	postResp := httptest.NewRecorder()
	body := new(bytes.Buffer)
	body.WriteString("{\"title\": \"test1\"}")
	req, _ := http.NewRequest("POST", "/v1/data", body)
	router.ServeHTTP(postResp, req)
	assert.Equal(t, 200, postResp.Code)

	getResp := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/data?title=test1", nil)
	router.ServeHTTP(getResp, req)
	assert.Equal(t, 200, getResp.Code)

	assert.Equal(t, postResp.Body.String(), getResp.Body.String())
}

func TestDataHandler_PostData_invalid(t *testing.T) {
	t.Parallel()

	testDB, _ := database.NewTestDatabase(t)

	h := &DataHandler{
		db: database.NewDataDB(testDB),
	}

	router := setupRouter(h)


	postResp := httptest.NewRecorder()
	body := new(bytes.Buffer)
	body.WriteString("{\"id\": \"test1\"}")
	req, _ := http.NewRequest("POST", "/v1/data", body)
	router.ServeHTTP(postResp, req)
	assert.Equal(t, 400, postResp.Code)
	assert.Equal(t, "{\"message\":\"Request field 'title' missing.\"}", postResp.Body.String())
}

func TestDataHandler_GetData_invalid(t *testing.T) {
	t.Parallel()

	testDB, _ := database.NewTestDatabase(t)

	h := &DataHandler{
		db: database.NewDataDB(testDB),
	}

	router := setupRouter(h)


	getResp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/data", nil)
	router.ServeHTTP(getResp, req)
	assert.Equal(t, 400, getResp.Code)
	assert.Equal(t, "{\"message\":\"Query param 'title' missing.\"}", getResp.Body.String())

	getResp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/data?title=nonexistent", nil)
	router.ServeHTTP(getResp, req)
	assert.Equal(t, 404, getResp.Code)
	assert.Equal(t, "{\"message\":\"No data record found with specified title\"}", getResp.Body.String())
}

func setupRouter(h *DataHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/v1/data", h.GetData)
	r.POST("/v1/data", h.PostData)
	return r
}
