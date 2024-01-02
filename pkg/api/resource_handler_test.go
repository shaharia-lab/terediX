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

func TestParseMetaDataEquals(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "Single key-value pair",
			input:    "key1=value1",
			expected: map[string]string{"key1": "value1"},
		},
		{
			name:     "Multiple key-value pairs",
			input:    "key1=value1,key2=value2",
			expected: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: map[string]string{},
		},
		{
			name:     "No value for key",
			input:    "key1=",
			expected: map[string]string{"key1": ""},
		},
		{
			name:     "No key for value",
			input:    "=value1",
			expected: map[string]string{"": "value1"},
		},
		{
			name:     "Extra equals signs",
			input:    "key1=value1=value2",
			expected: map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseMetaDataEquals(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
