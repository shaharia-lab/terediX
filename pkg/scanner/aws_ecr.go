// Package scanner scans targets
package scanner

import (
	"context"
	"fmt"

	ecrTypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/shaharia-lab/teredix/pkg"
	"github.com/shaharia-lab/teredix/pkg/config"
	"github.com/shaharia-lab/teredix/pkg/metrics"
	"github.com/shaharia-lab/teredix/pkg/resource"
	"github.com/shaharia-lab/teredix/pkg/storage"
	"github.com/shaharia-lab/teredix/pkg/util"
	"github.com/sirupsen/logrus"

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
	storage                storage.Storage
	logger                 *logrus.Logger
	schedule               string
	metrics                *metrics.Collector
}

// GetName return source name
func (a *AWSECR) GetName() string {
	return a.SourceName
}

// GetSchedule return schedule
func (a *AWSECR) GetSchedule() string {
	return a.schedule
}

// Setup AWS ECR source
func (a *AWSECR) Setup(name string, cfg config.Source, dependencies *Dependencies) error {
	a.SourceName = name
	a.ECRClient = ecr.NewFromConfig(buildAWSConfig(cfg))
	a.Region = cfg.Configuration["region"]
	a.AccountID = cfg.Configuration["account_id"]
	a.Fields = cfg.Fields
	a.ResourceTaggingService = resourcegroupstaggingapi.NewFromConfig(buildAWSConfig(cfg))
	a.storage = dependencies.GetStorage()
	a.logger = dependencies.GetLogger()
	a.metrics = dependencies.GetMetrics()

	a.logger.WithFields(logrus.Fields{
		"scanner_name": a.SourceName,
		"scanner_kind": a.GetKind(),
	}).Info("Scanner has been setup")

	return nil
}

// GetKind return resource kind
func (a *AWSECR) GetKind() string {
	return pkg.ResourceKindAWSECR
}

// Scan discover resource and send to resource channel
func (a *AWSECR) Scan(resourceChannel chan resource.Resource) error {
	nextVersion, err := a.storage.GetNextVersionForResource(a.SourceName, pkg.ResourceKindAWSECR)
	if err != nil {
		return fmt.Errorf("unable to get next version for resource: %w", err)
	}

	totalResourceDiscovered := 0

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
			a.logger.WithFields(logrus.Fields{
				"scanner_name": a.SourceName,
				"scanner_kind": a.GetKind(),
			}).WithError(err).Error("Unable to make api call to aws ecr endpoint")

			return fmt.Errorf("unable to make api call to aws ecr endpoint: %w", err)
		}

		// Loop through repositories and their tags
		for _, repository := range resp.Repositories {
			res := resource.NewResource(pkg.ResourceKindAWSECR, *repository.RepositoryName, *repository.RepositoryArn, a.SourceName, nextVersion)
			res.AddMetaData(a.getMetaData(repository))
			resourceChannel <- res

			totalResourceDiscovered++
		}

		// Check if there are more pages
		if resp.NextToken == nil {
			break
		}
		nextToken = *resp.NextToken
		pageNum++
	}

	a.logger.WithFields(logrus.Fields{
		"scanner_name":              a.SourceName,
		"scanner_kind":              a.GetKind(),
		"total_resource_discovered": totalResourceDiscovered,
	}).Info("scan completed")

	a.metrics.CollectTotalResourceDiscoveredByScanner(a.SourceName, a.GetKind(), float64(totalResourceDiscovered))
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
