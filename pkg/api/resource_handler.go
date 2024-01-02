package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
)

func GetAllResources(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		rResponse, err := getResources(s, r.URL.Query().Get("page"), r.URL.Query().Get("per_page"))
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

func getResources(s storage.Storage, page, perPage string) (resource.ListResponse, error) {
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

	// Use the Find method to retrieve resources
	resources, err := s.Find(filter)
	if err != nil {
		return resource.ListResponse{}, err
	}

	rResponse := resource.ListResponse{
		Page:      pageInt,
		PerPage:   perPageInt,
		HasMore:   true,
		Resources: []resource.Response{},
	}
	for _, re := range resources {
		res := resource.Response{
			Kind:       re.GetKind(),
			UUID:       re.GetUUID(),
			Name:       re.GetName(),
			ExternalID: re.GetExternalID(),
			Scanner:    re.GetScanner(),
			FetchedAt:  re.GetFetchedAt(),
			Version:    re.GetVersion(),
			MetaData:   map[string]string{},
		}

		rm := re.GetMetaData()
		if rm.Get() != nil {
			for _, m := range rm.Get() {
				res.MetaData[m.Key] = m.Value
			}
		}

		rResponse.Resources = append(rResponse.Resources, res)
	}
	return rResponse, nil
}
