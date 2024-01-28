package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource_AddRelation(t *testing.T) {
	res1 := NewResource("test", "test-resource-1", "test-id", "test-scanner", 1)
	res2 := NewResource("test", "test-resource-2", "test-id2", "test-scanner", 1)

	res1.AddRelation(res2)

	assert.Len(t, res1.relatedWith, 1)
	assert.Equal(t, "test", res1.relatedWith[0].GetKind())
	assert.Equal(t, "test-id2", res1.relatedWith[0].externalID)
}

func TestResource_AddMetaData(t *testing.T) {
	res1 := NewResource("test", "test-resource-1", "test-id", "test-scanner", 1)
	res1.AddMetaData(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	data := res1.GetMetaData()
	assert.Len(t, data.Get(), 2)
	assert.Equal(t, "value1", data.Find("key1").Value)
	assert.Equal(t, "value2", data.Find("key2").Value)
}

func TestResource_FindMetaValue(t *testing.T) {
	res1 := NewResource("test", "test-resource-1", "test-id", "test-scanner", 1)

	res1.AddMetaData(map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	value := res1.metaData.Find("key1")
	assert.Equal(t, "value1", value.Value)

	value = res1.metaData.Find("key3")
	assert.Nil(t, value)
}

func TestResource_ToAPIResponse(t *testing.T) {
	testCases := []struct {
		name     string
		resource Resource
		expected Response
	}{
		{
			name: "Test case 1: All fields are filled",
			resource: Resource{
				kind:       "testKind",
				uuid:       "testUUID",
				name:       "testName",
				externalID: "testExternalID",
				scanner:    "testScanner",
				version:    1,
				metaData:   MetaDataLists{data: []MetaData{{Key: "key1", Value: "value1"}}},
			},
			expected: Response{
				Kind:       "testKind",
				UUID:       "testUUID",
				Name:       "testName",
				ExternalID: "testExternalID",
				Scanner:    "testScanner",
				Version:    1,
				MetaData:   map[string]string{"key1": "value1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.resource.ToAPIResponse()
			assert.Equal(t, tc.expected, actual)
		})
	}
}
