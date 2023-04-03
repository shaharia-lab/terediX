// Package util represent few common functions
package util

import (
	"context"
	"fmt"
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
