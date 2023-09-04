// Package util represent few common functions
package util

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/google/uuid"
	"github.com/shaharia-lab/teredix/pkg/resource"
)

// GenerateUUID generate v4 uuid
func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}

// RetryWithExponentialBackoff retries a function with exponential backoff in case of errors
func RetryWithExponentialBackoff(fn func() error, maxRetries int, initialBackoffSeconds int) error {
	backoff := time.Duration(initialBackoffSeconds) * time.Second
	subsequentBackoff := 2

	for i := 0; ; i++ {
		fmt.Println("Retrying....")
		err := fn()
		if err == nil {
			return nil
		}

		if i == maxRetries {
			return fmt.Errorf("maximum number of retries exceeded: %w", err)
		}

		//log.Printf("Error occurred: %v. Retrying in %v", err, backoff)
		time.Sleep(backoff)

		backoff *= time.Duration(subsequentBackoff)
	}
}

// ResourceTaggingServiceClient construct resource tagging service
type ResourceTaggingServiceClient interface {
	GetResources(context.Context, *resourcegroupstaggingapi.GetResourcesInput, ...func(*resourcegroupstaggingapi.Options)) (*resourcegroupstaggingapi.GetResourcesOutput, error)
}

// GetAWSResourceTagByARN provides tags for any resource from ARN
func GetAWSResourceTagByARN(ctx context.Context, resourceTaggingService ResourceTaggingServiceClient, arn string) (map[string]string, error) {
	input := &resourcegroupstaggingapi.GetResourcesInput{
		ResourceARNList: []string{arn},
	}

	resp, err := resourceTaggingService.GetResources(ctx, input)
	if err != nil {
		return nil, err
	}

	tags := make(map[string]string)

	if len(resp.ResourceTagMappingList) == 0 {
		return tags, nil
	}

	for _, tag := range resp.ResourceTagMappingList[0].Tags {
		tags[*tag.Key] = *tag.Value
	}

	return tags, nil
}

// IsFieldExistsInConfig to check if specific field exists in config
func IsFieldExistsInConfig(value string, fields []string) bool {
	for _, v := range fields {
		if v == value {
			return true
		}
	}
	return false
}

// CheckKeysInMetaData Checks if all the keys in the given list exist in the MetaData of a Resource
// Returns a boolean indicating if all keys exist and a slice of missing keys
func CheckKeysInMetaData(resource resource.Resource, keys []string) (bool, []string) {
	var missingKeys []string

	data := resource.GetMetaData()
	for _, key := range keys {
		if data.Find(key) == nil {
			missingKeys = append(missingKeys, key)
		}
	}

	return len(missingKeys) == 0, missingKeys
}

// CheckIfMetaKeysExistsInResources Checks if all the keys in the given list exist in the metaData of all the Resources
func CheckIfMetaKeysExistsInResources(t *testing.T, res []resource.Resource, expectedMetaDataKeys []string) {
	for k, v := range res {
		exists, missingKeys := CheckKeysInMetaData(v, expectedMetaDataKeys)
		if !exists {
			t.Errorf("Metadata missing. Missing keys [%d]: %v", k, missingKeys)
		}
	}
}
