// Package api provides api handlers
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
)

// GetAllResources returns a handler function that retrieves all resources
func GetAllResources(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := buildResourceFilterFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resources, err := getResources(s, filter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(resources)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func buildResourceFilterFromRequest(r *http.Request) (storage.ResourceFilter, error) {
	page, perPage := getPageAndPerPage(r)

	filter := storage.ResourceFilter{
		PerPage:    perPage,
		Offset:     (page - 1) * perPage,
		Kind:       r.URL.Query().Get("kind"),
		Name:       r.URL.Query().Get("name"),
		ExternalID: r.URL.Query().Get("external_id"),
		UUID:       r.URL.Query().Get("uuid"),
	}

	if metaDataEq := r.URL.Query().Get("meta_data_eq"); metaDataEq != "" {
		filter.MetaDataEquals = parseMetaDataEquals(metaDataEq)
	}

	return filter, nil
}

func getPageAndPerPage(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	if pageStr == "" {
		pageStr = "1"
	}

	if perPageStr == "" {
		perPageStr = "200"
	}

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)

	if perPage > 200 {
		perPage = 200
	}

	return page, perPage
}

func getResources(s storage.Storage, filter storage.ResourceFilter) (resource.ListResponse, error) {
	resources, err := s.Find(filter)
	if err != nil {
		return resource.ListResponse{}, err
	}

	rResponse := resource.ListResponse{
		Page:      filter.Offset + 1,
		PerPage:   filter.PerPage,
		HasMore:   true,
		Resources: []resource.Response{},
	}
	for _, re := range resources {
		rResponse.Resources = append(rResponse.Resources, re.ToAPIResponse())
	}
	return rResponse, nil
}

func parseMetaDataEquals(metaDataEq string) map[string]string {
	metaData := make(map[string]string)

	// Split the string by semicolon to get key-value pairs
	pairs := strings.Split(metaDataEq, ",")

	for _, pair := range pairs {
		// Split each pair by equals sign to get the key and value
		kv := strings.Split(pair, "=")

		// Check that the pair was split into exactly two parts
		if len(kv) == 2 {
			// Add the key-value pair to the map
			metaData[kv[0]] = kv[1]
		}
	}

	return metaData
}
