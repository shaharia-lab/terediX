package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource_AddRelation(t *testing.T) {
	res1 := NewResource("test", "test-resource-1", "test-id", "test-scanner", "1.0.0")
	res2 := NewResource("test", "test-resource-2", "test-id2", "test-scanner", "1.0.0")

	res1.AddRelation(res2)

	assert.Len(t, res1.RelatedWith, 1)
	assert.Equal(t, "test", res1.RelatedWith[0].Kind)
	assert.Equal(t, "test-id2", res1.RelatedWith[0].ExternalID)
}

func TestResource_AddMetaData(t *testing.T) {
	res1 := NewResource("test", "test-resource-1", "test-id", "test-scanner", "1.0.0")

	res1.AddMetaData("key1", "value1")
	res1.AddMetaData("key2", "value2")

	assert.Len(t, res1.MetaData, 2)
	assert.Equal(t, "key1", res1.MetaData[0].Key)
	assert.Equal(t, "value1", res1.MetaData[0].Value)
	assert.Equal(t, "key2", res1.MetaData[1].Key)
	assert.Equal(t, "value2", res1.MetaData[1].Value)
}

func TestResource_FindMetaValue(t *testing.T) {
	res1 := NewResource("test", "test-resource-1", "test-id", "test-scanner", "1.0.0")

	res1.AddMetaData("key1", "value1")
	res1.AddMetaData("key2", "value2")

	value := res1.FindMetaValue("key1")
	assert.Equal(t, "value1", value)

	value = res1.FindMetaValue("key3")
	assert.Equal(t, "", value)
}
