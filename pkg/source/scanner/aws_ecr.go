// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"
	"strings"

	ecrTypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

const (
	ecrFieldRepositoryName = "repositoryName"
	ecrFieldArn            = "repositoryArn"
	ecrFieldRegistryId     = "registryId"
	ecrFieldRepositoryUri  = "repositoryUri"
	ecrFieldTags           = "tags"
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
	Fields                 []string
}

// NewAWSECR construct AWS ECR source
func NewAWSECR(sourceName string, region string, accountID string, ecrClient EcrClient, resourceTaggingService util.ResourceTaggingServiceClient, fields []string) *AWSECR {
	return &AWSECR{
		SourceName:             sourceName,
		ECRClient:              ecrClient,
		Region:                 region,
		AccountID:              accountID,
		ResourceTaggingService: resourceTaggingService,
		Fields:                 fields,
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
			res, err := a.mapToResource(repository)
			if err != nil {
				return err
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

func (a *AWSECR) mapToResource(repository ecrTypes.Repository) (resource.Resource, error) {
	var repoMeta []resource.MetaData
	for _, mapper := range a.fieldMapper(repository) {
		if util.IsFieldExistsInConfig(mapper.field, a.Fields) || strings.Contains(mapper.field, "tag_") {
			val := mapper.value()
			if val != "" {
				repoMeta = append(repoMeta, resource.MetaData{Key: mapper.field, Value: val})
			}
		}
	}

	return resource.Resource{
		Name:       *repository.RepositoryName,
		Kind:       pkg.ResourceKindAWSECR,
		UUID:       util.GenerateUUID(),
		ExternalID: *repository.RepositoryArn,
		MetaData:   repoMeta,
	}, nil
}

func (a *AWSECR) fieldMapper(repository ecrTypes.Repository) []MetaDataMapper {
	var fieldMapper []MetaDataMapper

	mappings := map[string]func() string{
		ecrFieldRepositoryName: func() string { return stringValueOrDefault(*repository.RepositoryName) },
		ecrFieldArn:            func() string { return stringValueOrDefault(*repository.RepositoryArn) },
		ecrFieldRegistryId:     func() string { return stringValueOrDefault(*repository.RegistryId) },
		ecrFieldRepositoryUri:  func() string { return stringValueOrDefault(*repository.RepositoryUri) },
	}

	for field, fn := range mappings {
		fieldMapper = append(fieldMapper, MetaDataMapper{field: field, value: fn})
	}

	if util.IsFieldExistsInConfig(fieldTags, a.Fields) {
		tags, err := util.GetAWSResourceTagByARN(context.Background(), a.ResourceTaggingService, *repository.RepositoryArn)
		if err != nil {
			return fieldMapper
		}
		for key, val := range tags {
			fieldMapper = append(fieldMapper, MetaDataMapper{
				field: fmt.Sprintf("tag_%s", key),
				value: func() string { return stringValueOrDefault(val) },
			})
		}
	}

	return fieldMapper
}
