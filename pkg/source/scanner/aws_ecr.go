// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// EcrClient build aws client
type EcrClient interface {
	DescribeRepositories(context.Context, *ecr.DescribeRepositoriesInput, ...func(options *ecr.Options)) (*ecr.DescribeRepositoriesOutput, error)
	GetRepositoryPolicy(ctx context.Context, params *ecr.GetRepositoryPolicyInput, optFns ...func(*ecr.Options)) (*ecr.GetRepositoryPolicyOutput, error)
	DescribeImages(context.Context, *ecr.DescribeImagesInput, ...func(*ecr.Options)) (*ecr.DescribeImagesOutput, error)
}

// AWSECR AWS ECR source
type AWSECR struct {
	SourceName             string
	ECRClient              EcrClient
	Region                 string
	AccountID              string
	ResourceTaggingService util.ResourceTaggingServiceClient
}

// NewAWSECR construct AWS ECR source
func NewAWSECR(sourceName string, region string, accountID string, ecrClient EcrClient, resourceTaggingService util.ResourceTaggingServiceClient) *AWSECR {
	return &AWSECR{
		SourceName:             sourceName,
		ECRClient:              ecrClient,
		Region:                 region,
		AccountID:              accountID,
		ResourceTaggingService: resourceTaggingService,
	}
}

// Scan discover resource and send to resource channel
// Scan discover resource and send to resource channel
func (a *AWSECR) Scan(resourceChannel chan resource.Resource) error {
	// Set initial values for pagination
	pageNum := 0
	nextToken := ""

	// Loop through pages of ECR repositories
	for {
		// Describe repositories for current page
		params := &ecr.DescribeRepositoriesInput{
			MaxResults: aws.Int32(int32(perPage)),
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := a.ECRClient.DescribeRepositories(context.TODO(), params)
		if err != nil {
			return err
		}

		// Loop through repositories and their tags
		for _, repository := range resp.Repositories {
			res := resource.Resource{
				Name:       *repository.RepositoryName,
				Kind:       pkg.ResourceKindAWSECR,
				UUID:       util.GenerateUUID(),
				ExternalID: *repository.RepositoryArn,
				MetaData: []resource.MetaData{
					{
						Key:   "AWS-ECR-Repository-Name",
						Value: *repository.RepositoryName,
					},
					{
						Key:   "AWS-ECR-Repository-Arn",
						Value: *repository.RepositoryArn,
					},
					{
						Key:   "AWS-ECR-Registry-Id",
						Value: *repository.RegistryId,
					},
					{
						Key:   "AWS-ECR-Repository-URI",
						Value: *repository.RepositoryUri,
					},
					{
						Key:   pkg.MetaKeyScannerLabel,
						Value: a.SourceName,
					},
				},
			}

			tags, err := util.GetAWSResourceTagByARN(context.Background(), a.ResourceTaggingService, *repository.RepositoryArn)
			if err != nil {
				return err
			}

			for tagKey, tagValue := range tags {
				res.MetaData = append(res.MetaData, resource.MetaData{
					Key:   fmt.Sprintf("AWS-ECR-%s", tagKey),
					Value: tagValue,
				})
			}

			resourceChannel <- res
		}

		// Check if there are more pages
		if resp.NextToken == nil {
			break
		}
		nextToken = *resp.NextToken
		pageNum++
	}

	return nil
}
