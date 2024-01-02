package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	s.setupAPIServer("8080")

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

	r1 := resource.NewResource("vm", "test", "test", "test", 1)
	st.On("Find", mock.Anything).Return([]resource.Resource{r1}, nil)

	s := NewServer(logger, st)
	s.setupAPIServer("8080")

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

func TestServer_SetupPromMetricsServer(t *testing.T) {
	logger := logrus.New()
	s := NewServer(logger, nil)

	s.setupPromMetricsServer()

	assert.NotNil(t, s.promMetricsServer)
	assert.Equal(t, ":2112", s.promMetricsServer.Addr)
}

func TestServer_PromMetricsEndpoint(t *testing.T) {
	logger := logrus.New()
	s := NewServer(logger, nil)
	s.setupPromMetricsServer()

	ts := httptest.NewServer(s.promMetricsServer.Handler)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/metrics")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestServer_StartAndShutdownServer(t *testing.T) {
	// Arrange
	logger := logrus.New()
	s := NewServer(logger, nil)
	s.setupAPIServer("0")

	// Act
	// Start the server in a separate goroutine so it doesn't block
	go s.startServer(s.apiServer, "API server")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Send a request to the server
	ts := httptest.NewServer(s.apiServer.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Act
	// Shutdown the server
	err = s.shutdownServer(ctx, s.apiServer, "API server")

	// Assert
	assert.NoError(t, err)

	// Wait for the server to shutdown
	time.Sleep(100 * time.Millisecond)

	// Check if the server is still running
	if s.apiServer != nil {
		t.Errorf("Server is still running after shutdown")
	}
}
