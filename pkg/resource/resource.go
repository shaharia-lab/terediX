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

type MetaDataLists struct {
	MetaData []MetaData
}

func (ml *MetaDataLists) IsExists(key string) bool {
	for _, metaData := range ml.MetaData {
		if metaData.Key == key {
			return true
		}
	}

	return false
}

func (ml *MetaDataLists) Add(key, value string) {
	if ml.IsExists(key) {
		return
	}

	ml.MetaData = append(ml.MetaData, MetaData{Key: value})
}

func (ml *MetaDataLists) AddMap(metaMap map[string]string) {
	for k, v := range metaMap {
		ml.Add(k, v)
	}
}

func (ml *MetaDataLists) Find(key string) *MetaData {
	if !ml.IsExists(key) {
		return nil
	}

	for _, v := range ml.MetaData {
		if v.Key == key {
			return &v
		}
	}

	return nil
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

func (r *Resource) SetUUID(uuid string) {
	r.UUID = uuid
}

func (r *Resource) GetScanner() string {
	return r.scanner
}

func (r *Resource) GetVersion() string {
	return r.version
}

func (r *Resource) GetFetchedAt() time.Time {
	return r.fetchedAt
}

func (r *Resource) GetMetaData() []MetaData {
	return r.MetaData
}

func (r *Resource) GetKind() string {
	return r.Kind
}

func (r *Resource) GetExternalID() string {
	return r.ExternalID
}

func (r *Resource) GetUUID() string {
	return r.UUID
}

func (r *Resource) GetName() string {
	return r.Name
}

func (r *Resource) GetRelatedResources() []Resource {
	return r.RelatedWith
}

// AddRelation build relation between resources
func (r *Resource) AddRelation(relatedResource Resource) {
	r.RelatedWith = append(r.RelatedWith, relatedResource)
}

// AddMetaDataMultiple adds or updates metadata for each resource
func (r *Resource) AddMetaDataMultiple(metaMap map[string]string) {
	for k, v := range metaMap {
		r.MetaData = append(r.MetaData, MetaData{
			Key:   k,
			Value: v,
		})
	}
}

func (r *Resource) FindMetaValue(key string) string {
	for _, v := range r.MetaData {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}
