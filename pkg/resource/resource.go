// Package resource represents resources
package resource

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

// AddMetaData adds meta data for each resource
func (r *Resource) AddMetaData(key, value string) {
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
