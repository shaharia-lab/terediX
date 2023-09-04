package resource

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

	ml.MetaData = append(ml.MetaData, MetaData{Key: key, Value: value})
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
