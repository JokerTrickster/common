package aws

import (
	"context"
	"sync"

	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type AWSService struct {
	region        string
	awsConfig     aws.Config
	initOnce      sync.Once
	ssmClient     *ssm.Client
	sesClient     *sesv2.Client
	s3Client      *s3.Client
	uploader      *manager.Uploader
	downloader    *manager.Downloader
	presignClient *s3.PresignClient
}

var instance *AWSService
var once sync.Once

// GetAWSService returns a singleton AWSService instance
func GetAWSService(region string) *AWSService {
	once.Do(func() {
		instance = &AWSService{region: region}
	})
	return instance
}

// Initialize loads the AWS configuration once
func (a *AWSService) Initialize() error {
	var err error
	a.initOnce.Do(func() {
		a.awsConfig, err = AwsConfig.LoadDefaultConfig(context.TODO(), AwsConfig.WithRegion(a.region))
	})
	return err
}

// GetSSMClient lazily initializes and returns the SSM client
func (a *AWSService) GetSSMClient() (*ssm.Client, error) {
	if a.ssmClient == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
		a.ssmClient = ssm.NewFromConfig(a.awsConfig)
	}
	return a.ssmClient, nil
}

// GetSESClient lazily initializes and returns the SES client
func (a *AWSService) GetSESClient() (*sesv2.Client, error) {
	if a.sesClient == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
		a.sesClient = sesv2.NewFromConfig(a.awsConfig)
	}
	return a.sesClient, nil
}

// GetS3Client lazily initializes and returns the S3 client
func (a *AWSService) GetS3Client() (*s3.Client, error) {
	if a.s3Client == nil {
		if err := a.Initialize(); err != nil {
			return nil, err
		}
		a.s3Client = s3.NewFromConfig(a.awsConfig)
	}
	return a.s3Client, nil
}

// GetUploader lazily initializes and returns the S3 uploader
func (a *AWSService) GetUploader() (*manager.Uploader, error) {
	if a.uploader == nil {
		s3Client, err := a.GetS3Client()
		if err != nil {
			return nil, err
		}
		a.uploader = manager.NewUploader(s3Client)
	}
	return a.uploader, nil
}

// GetDownloader lazily initializes and returns the S3 downloader
func (a *AWSService) GetDownloader() (*manager.Downloader, error) {
	if a.downloader == nil {
		s3Client, err := a.GetS3Client()
		if err != nil {
			return nil, err
		}
		a.downloader = manager.NewDownloader(s3Client)
	}
	return a.downloader, nil
}

// GetPresignClient lazily initializes and returns the S3 presign client
func (a *AWSService) GetPresignClient() (*s3.PresignClient, error) {
	if a.presignClient == nil {
		s3Client, err := a.GetS3Client()
		if err != nil {
			return nil, err
		}
		a.presignClient = s3.NewPresignClient(s3Client)
	}
	return a.presignClient, nil
}
