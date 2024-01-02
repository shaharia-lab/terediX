// Package api provides api handlers
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
)

// GetAllResources returns a handler function that retrieves all resources
func GetAllResources(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		perPage := r.URL.Query().Get("per_page")

		if page == "" {
			page = "1"
		}

		if perPage == "" {
			perPage = "200"
		}

		// Convert query parameters to integers
		pageInt, _ := strconv.Atoi(page)
		perPageInt, _ := strconv.Atoi(perPage)

		// Ensure perPage does not exceed 200
		if perPageInt > 200 {
			perPageInt = 200
		}

		// Create a ResourceFilter
		filter := storage.ResourceFilter{PerPage: perPageInt, Offset: (pageInt - 1) * perPageInt}

		if kind := r.URL.Query().Get("kind"); kind != "" {
			filter.Kind = kind
		}

		if name := r.URL.Query().Get("name"); name != "" {
			filter.Name = name
		}

		if externalID := r.URL.Query().Get("external_id"); externalID != "" {
			filter.ExternalID = externalID
		}

		if uuid := r.URL.Query().Get("uuid"); uuid != "" {
			filter.UUID = uuid
		}

		if metaDataEq := r.URL.Query().Get("meta_data_eq"); metaDataEq != "" {
			filter.MetaDataEquals = parseMetaDataEquals(metaDataEq)
		}

		for k, v := range filter.MetaDataEquals {
			log.Print(k, v)
		}

		// Parse query parameters
		rResponse, err := getResources(s, filter, pageInt, perPageInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the response to JSON
		jsonResponse, err := json.Marshal(rResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func getResources(s storage.Storage, filter storage.ResourceFilter, page, perPage int) (resource.ListResponse, error) {
	// Use the Find method to retrieve resources
	resources, err := s.Find(filter)
	if err != nil {
		return resource.ListResponse{}, err
	}

	rResponse := resource.ListResponse{
		Page:      page,
		PerPage:   perPage,
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
