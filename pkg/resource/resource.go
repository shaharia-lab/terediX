package resource

type Resource struct {
	Kind        string
	UUID        string
	Name        string
	ExternalID  string
	RelatedWith []Resource
	MetaData    []MetaData
}

type MetaData struct {
	Key   string
	Value string
}

type ChildOf struct {
	Kind         string
	ResourceUUID string
}

func NewResource(kind, uuid, name, externalID string, scanner string) Resource {
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

func (r *Resource) AddRelation(relatedResource Resource) {
	r.RelatedWith = append(r.RelatedWith, relatedResource)
}

func (r *Resource) AddMetaData(key, value string) {
	r.MetaData = append(r.MetaData, MetaData{
		Key:   key,
		Value: value,
	})
}

func (r *Resource) FindMetaValue(key string) string {
	for _, v := range r.MetaData {
		if v.Key == key {
			return v.Value
		}
	}
	return ""
}
