package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_SetupAPIServer(t *testing.T) {
	logger := logrus.New()
	st := new(storage.Mock)

	s := NewServer(logger, st)

	s.setupAPIServer()

	assert.NotNil(t, s.apiServer)
	assert.Equal(t, ":8080", s.apiServer.Addr)
}

func TestServer_GetResources(t *testing.T) {
	logger := logrus.New()
	st := new(storage.Mock)
	st.On("Find", mock.Anything).Return([]resource.Resource{}, nil)

	s := NewServer(logger, st)

	// Test with valid page and per_page parameters
	response, err := s.getResources("1", "10")
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PerPage)

	// Test with empty page and per_page parameters
	response, err = s.getResources("", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 200, response.PerPage)

	// Test with per_page parameter greater than 200
	response, err = s.getResources("1", "300")
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 200, response.PerPage)
}

func TestServer_APIEndpoint(t *testing.T) {
	logger := logrus.New()
	st := new(storage.Mock)
	st.On("Find", mock.Anything).Return([]resource.Resource{}, nil)

	s := NewServer(logger, st)
	s.setupAPIServer()

	r := chi.NewRouter()
	r.Mount("/", s.apiServer.Handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/v1/resources?page=1&per_page=10")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/api/v1/resources?page=1&per_page=300")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/api/v1/resources")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
