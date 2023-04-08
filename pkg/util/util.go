// Package util represent few common functions
package util

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"

	"github.com/google/uuid"
)

// GenerateUUID generate v4 UUID
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

// IsExist check whether the given value exists in a slice
func IsExist(what interface{}, in interface{}) bool {
	s := reflect.ValueOf(in)

	if s.Kind() != reflect.Slice {
		panic("IsExist: Second argument must be a slice")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Kind() != reflect.TypeOf(what).Kind() {
			continue
		}

		switch s.Index(i).Kind() {
		case reflect.String:
			if s.Index(i).String() == reflect.ValueOf(what).String() {
				return true
			}
		default:
			if reflect.DeepEqual(what, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}
