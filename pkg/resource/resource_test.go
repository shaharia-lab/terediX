package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource_AddRelation(t *testing.T) {
	res1 := NewResource("test", "1", "test-resource-1", "test-id", "test-scanner")
	res2 := NewResource("test", "2", "test-resource-2", "test-id", "test-scanner")

	res1.AddRelation(res2)

	assert.Len(t, res1.RelatedWith, 1)
	assert.Equal(t, "test", res1.RelatedWith[0].Kind)
	assert.Equal(t, "2", res1.RelatedWith[0].UUID)
}

func TestResource_AddMetaData(t *testing.T) {
	res1 := NewResource("test", "1", "test-resource-1", "test-id", "test-scanner")

	res1.AddMetaData("key1", "value1")
	res1.AddMetaData("key2", "value2")

	assert.Len(t, res1.MetaData, 3)
	assert.Equal(t, "Scanner", res1.MetaData[0].Key)
	assert.Equal(t, "test-scanner", res1.MetaData[0].Value)
	assert.Equal(t, "key1", res1.MetaData[1].Key)
	assert.Equal(t, "value1", res1.MetaData[1].Value)
	assert.Equal(t, "key2", res1.MetaData[2].Key)
	assert.Equal(t, "value2", res1.MetaData[2].Value)
}

func TestResource_FindMetaValue(t *testing.T) {
	res1 := NewResource("test", "1", "test-resource-1", "test-id", "test-scanner")

	res1.AddMetaData("key1", "value1")
	res1.AddMetaData("key2", "value2")

	value := res1.FindMetaValue("key1")
	assert.Equal(t, "value1", value)

	value = res1.FindMetaValue("key3")
	assert.Equal(t, "", value)
}
