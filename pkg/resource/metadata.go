// Package resource represents resources
package resource

// MetaData resource metadata
type MetaData struct {
	Key   string
	Value string
}

// MetaDataLists list of metadata
type MetaDataLists struct {
	data []MetaData
}

// IsExists checks if metadata exists
func (ml *MetaDataLists) IsExists(key string) bool {
	for _, metaData := range ml.data {
		if metaData.Key == key {
			return true
		}
	}

	return false
}

// Get checks if metadata exists
func (ml *MetaDataLists) Get() []MetaData {
	return ml.data
}

// Add adds metadata
func (ml *MetaDataLists) Add(key, value string) {
	if ml.IsExists(key) {
		return
	}

	ml.data = append(ml.data, MetaData{Key: key, Value: value})
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

	for _, v := range ml.data {
		if v.Key == key {
			return &v
		}
	}

	return nil
}

// FindMissingKeys returns missing keys
func (ml *MetaDataLists) FindMissingKeys(keys []string) []string {
	var missingKeys []string

	for _, key := range keys {
		if ml.Find(key) == nil {
			missingKeys = append(missingKeys, key)
		}
	}

	return missingKeys
}
