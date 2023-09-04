// Package resource represents resources
package resource

import (
	"time"

	"github.com/google/uuid"
)

// Resource represent resource
type Resource struct {
	kind        string
	uuid        string
	name        string
	externalID  string
	relatedWith []Resource
	metaData    MetaDataLists
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
func (r *Resource) GetMetaData() MetaDataLists {
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
		r.metaData.Add(k, v)
	}
}
