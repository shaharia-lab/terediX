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

// MetaDataLists list of metadata
type MetaDataLists struct {
	MetaData []MetaData
}

// IsExists checks if metadata exists
func (ml *MetaDataLists) IsExists(key string) bool {
	for _, metaData := range ml.MetaData {
		if metaData.Key == key {
			return true
		}
	}

	return false
}

// Add adds metadata
func (ml *MetaDataLists) Add(key, value string) {
	if ml.IsExists(key) {
		return
	}

	ml.MetaData = append(ml.MetaData, MetaData{Key: value})
}

// AddMap adds metadata from map
func (ml *MetaDataLists) AddMap(metaMap map[string]string) {
	for k, v := range metaMap {
		ml.Add(k, v)
	}
}

// Find returns metadata by key
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
	kind        string
	uuid        string
	name        string
	externalID  string
	relatedWith []Resource
	metaData    []MetaData
	scanner     string
	fetchedAt   time.Time
	version     string
}

// NewResource instantiate new resource
func NewResource(kind, name, externalID, scannerName, version string) Resource {
	return Resource{
		kind:       kind,
		name:       name,
		uuid:       uuid.New().String(),
		externalID: externalID,
		version:    version,
		scanner:    scannerName,
		fetchedAt:  time.Now().UTC(),
	}
}

// SetUUID sets resource UUID
func (r *Resource) SetUUID(uuid string) {
	r.uuid = uuid
}

// GetScanner returns resource scanner
func (r *Resource) GetScanner() string {
	return r.scanner
}

// GetVersion returns resource version
func (r *Resource) GetVersion() string {
	return r.version
}

// GetFetchedAt returns resource fetched at
func (r *Resource) GetFetchedAt() time.Time {
	return r.fetchedAt
}

// GetMetaData returns resource metadata
func (r *Resource) GetMetaData() []MetaData {
	return r.metaData
}

// GetKind returns resource kind
func (r *Resource) GetKind() string {
	return r.kind
}

// GetExternalID returns resource external ID
func (r *Resource) GetExternalID() string {
	return r.externalID
}

// GetUUID returns resource UUID
func (r *Resource) GetUUID() string {
	return r.uuid
}

// GetName returns resource name
func (r *Resource) GetName() string {
	return r.name
}

// GetRelatedResources returns related resources
func (r *Resource) GetRelatedResources() []Resource {
	return r.relatedWith
}

// AddRelation build relation between resources
func (r *Resource) AddRelation(relatedResource Resource) {
	r.relatedWith = append(r.relatedWith, relatedResource)
}

// AddMetaData adds or updates metadata for a resource
func (r *Resource) AddMetaData(metaDataMap map[string]string) {
	for k, v := range metaDataMap {
		r.metaData = append(r.metaData, MetaData{
			Key:   k,
			Value: v,
		})
	}
}

// AddMetaDataMultiple adds or updates metadata for each resource
func (r *Resource) AddMetaDataMultiple(metaMap map[string]string) {
	for k, v := range metaMap {
		r.metaData = append(r.metaData, MetaData{
			Key:   k,
			Value: v,
		})
	}
}

// FindMetaValue finds the value of a metadata key
func (r *Resource) FindMetaValue(key string) string {
	for _, v := range r.metaData {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}
