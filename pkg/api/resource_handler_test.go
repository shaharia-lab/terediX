package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllResources(t *testing.T) {
	st := new(storage.Mock)
	st.On("Find", mock.Anything).Return([]resource.Resource{}, nil)

	handler := GetAllResources(st)

	req, err := http.NewRequest("GET", "/resources?page=1&per_page=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	req, err = http.NewRequest("GET", "/resources", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	req, err = http.NewRequest("GET", "/resources?page=1&per_page=300", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetResources(t *testing.T) {
	st := new(storage.Mock)
	st.On("Find", mock.Anything).Return([]resource.Resource{}, nil)

	// Test with valid page and per_page parameters
	response, err := getResources(st, "1", "10")
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PerPage)

	// Test with empty page and per_page parameters
	response, err = getResources(st, "", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 200, response.PerPage)

	// Test with per_page parameter greater than 200
	response, err = getResources(st, "1", "300")
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 200, response.PerPage)
}
