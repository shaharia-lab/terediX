// Package resource represents resources
package resource

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Resource represent resource
type Resource struct {
	Kind        string
	UUID        string
	Name        string
	ExternalID  string
	RelatedWith []Resource
	MetaData    []MetaData
}

// MetaData resource metadata
type MetaData struct {
	Key   string
	Value string
}

func NewResource(kind, name, externalID, scannerName, version string) Resource {
	return Resource{
		Kind:       kind,
		Name:       name,
		UUID:       uuid.New().String(),
		ExternalID: externalID,
		MetaData: []MetaData{
			{
				Key:   "_scanner",
				Value: scannerName,
			},
			{
				Key:   "_fetched_at",
				Value: strconv.FormatInt(time.Now().UTC().Unix(), 10),
			},
			{
				Key:   "_version",
				Value: version,
			},
		},
	}
}

// NewResourceV1 construct new resource
func NewResourceV1(kind, uuid, name, externalID string, scanner string) Resource {
	return Resource{
		Kind:       kind,
		UUID:       uuid,
		Name:       name,
		ExternalID: externalID,
		MetaData: []MetaData{
			{
				Key:   "Scanner",
				Value: scanner,
			},
		},
		RelatedWith: []Resource{},
	}
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
