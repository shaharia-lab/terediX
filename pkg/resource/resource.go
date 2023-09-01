// Package resource represents resources
package resource

import (
	"time"

	"github.com/google/uuid"
)

// MetaData resource metadata
type MetaData struct {
	Key   string
	Value string
}

// Resource represent resource
type Resource struct {
	Kind        string
	UUID        string
	Name        string
	ExternalID  string
	RelatedWith []Resource
	MetaData    []MetaData
	scanner     string
	fetchedAt   time.Time
	version     string
}

// NewResource instantiate new resource
func NewResource(kind, name, externalID, scannerName, version string) Resource {
	return Resource{
		Kind:       kind,
		Name:       name,
		UUID:       uuid.New().String(),
		ExternalID: externalID,
		version:    version,
		scanner:    scannerName,
		fetchedAt:  time.Now().UTC(),
	}
}

func (r *Resource) GetScanner() string {
	return r.scanner
}

// AddRelation build relation between resources
func (r *Resource) AddRelation(relatedResource Resource) {
	r.RelatedWith = append(r.RelatedWith, relatedResource)
}

// AddMetaData adds or updates meta data for each resource
func (r *Resource) AddMetaData(key, value string) {
	// Loop through existing MetaData to see if key already exists
	for i, metaData := range r.MetaData {
		if metaData.Key == key {
			// If key exists, update its value and return
			r.MetaData[i].Value = value
			return
		}
	}

	// If the key doesn't exist, add it to MetaData
	r.MetaData = append(r.MetaData, MetaData{
		Key:   key,
		Value: value,
	})
}

// FindMetaValue finds meta value by key
func (r *Resource) FindMetaValue(key string) string {
	for _, v := range r.MetaData {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}

func mapToMetaData(metaMap map[string]string) []MetaData {
	var metaData []MetaData
	for k, v := range metaMap {
		metaData = append(metaData, MetaData{
			Key:   k,
			Value: v,
		})
	}
	return metaData
}
