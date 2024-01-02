package api

import (
	"testing"

	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllResources(t *testing.T) {
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
