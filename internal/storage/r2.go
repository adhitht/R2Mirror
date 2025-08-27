package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"	
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	
	cfg "github.com/adhitht/R2Mirror/internal/config"
)

type R2Client struct {
	s3Client *s3.Client
	uploader *manager.Uploader
	config   *cfg.Config
}

func NewR2Client(config *cfg.Config) (*R2Client, error) {
	creds, err := cfg.GetStorageCredentials()
	if err != nil {
		return nil, err
	}

	configOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(config.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			creds.AccessKeyID, 
			creds.SecretAccessKey, 
			"",
		)),
	}

	if creds.EndpointURL != "" {
		fmt.Printf("🌐 Using R2 endpoint: %s\n", creds.EndpointURL)
		configOptions = append(configOptions, awsconfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               creds.EndpointURL,
					SigningRegion:     config.Region,
					HostnameImmutable: true,
				}, nil
			}),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to load R2 config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	uploader := manager.NewUploader(s3Client)

	return &R2Client{
		s3Client: s3Client,
		uploader: uploader,
		config:   config,
	}, nil
}

func (r *R2Client) UploadFile(ctx context.Context, bucket, key string, body io.Reader, contentType string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	_, err := r.uploader.Upload(ctx, input)
	return err
}

func (r *R2Client) UploadHTML(ctx context.Context, bucket, key, content string) error {
	return r.UploadFile(ctx, bucket, key, strings.NewReader(content), "text/html")
}

func (r *R2Client) GetPublicURL(bucket, key string) string {
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", bucket, key)
}

func (r *R2Client) Close() error {
	return nil
}