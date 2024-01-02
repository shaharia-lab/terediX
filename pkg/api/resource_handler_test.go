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

	testCases := []struct {
		name     string
		page     string
		perPage  string
		expected int
	}{
		{"valid page and per_page parameters", "1", "10", http.StatusOK},
		{"empty page and per_page parameters", "", "", http.StatusOK},
		{"per_page parameter greater than 200", "1", "300", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/resources?page="+tc.page+"&per_page="+tc.perPage, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expected, rr.Code)
		})
	}
}
