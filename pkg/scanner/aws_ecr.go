// Package scanner scans targets
package scanner

import (
	"context"

	ecrTypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/go-co-op/gocron"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/sirupsen/logrus"

	"github.com/shaharia-lab/teredix/pkg/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

const (
	ecrFieldRepositoryName = "repositoryName"
	ecrFieldArn            = "repositoryArn"
	ecrFieldRegistryID     = "registryID"
	ecrFieldRepositoryURI  = "repositoryURI"
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
	Schedule               config.Schedule
	scheduler              *gocron.Scheduler
	storage                storage.Storage
	logger                 *logrus.Logger
}

func (a *AWSECR) setECRClient(ecrClient EcrClient) {
	a.ECRClient = ecrClient
}

func (a *AWSECR) setResourceTaggingService(resourceTaggingService util.ResourceTaggingServiceClient) {
	a.ResourceTaggingService = resourceTaggingService
}

// Build AWS ECR source
func (a *AWSECR) Build(sourceKey string, cfg config.Source, storage storage.Storage, scheduler *gocron.Scheduler, logger *logrus.Logger) Scanner {
	a.SourceName = sourceKey
	a.ECRClient = ecr.NewFromConfig(BuildAWSConfig(cfg))
	a.Region = cfg.Configuration["region"]
	a.AccountID = cfg.Configuration["account_id"]
	a.Fields = cfg.Fields
	a.Schedule = cfg.Schedule
	a.scheduler = scheduler
	a.storage = storage
	a.logger = logger

	return a
}

// GetKind return resource kind
func (a *AWSECR) GetKind() string {
	return pkg.ResourceKindAWSECR
}

// Scan discover resource and send to resource channel
// Scan discover resource and send to resource channel
func (a *AWSECR) Scan(resourceChannel chan resource.Resource) error {
	// Set initial values for pagination
	pageNum := 0
	nextToken := ""

	nextVersion, err := a.storage.GetNextVersionForResource(a.SourceName, pkg.ResourceKindAWSEC2)
	if err != nil {
		return err
	}

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
			res := resource.NewResource(pkg.ResourceKindAWSECR, *repository.RepositoryName, *repository.RepositoryArn, a.SourceName, nextVersion)
			res.AddMetaData(a.getMetaData(repository))
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

func (a *AWSECR) getMetaData(repository ecrTypes.Repository) map[string]string {
	mappings := map[string]func() string{
		ecrFieldRepositoryName: func() string { return stringValueOrDefault(*repository.RepositoryName) },
		ecrFieldArn:            func() string { return stringValueOrDefault(*repository.RepositoryArn) },
		ecrFieldRegistryID:     func() string { return stringValueOrDefault(*repository.RegistryId) },
		ecrFieldRepositoryURI:  func() string { return stringValueOrDefault(*repository.RepositoryUri) },
	}

	getTags := func() []ResourceTag {
		tags, err := util.GetAWSResourceTagByARN(context.Background(), a.ResourceTaggingService, *repository.RepositoryArn)
		if err != nil {
			return []ResourceTag{}
		}

		var tt []ResourceTag
		for key, val := range tags {
			tt = append(tt, ResourceTag{
				Key:   key,
				Value: val,
			})
		}

		return tt
	}

	return NewFieldMapper(mappings, getTags, a.Fields).getResourceMetaData()
}
