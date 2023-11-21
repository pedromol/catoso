package storage

import (
	"bytes"
	"context"

	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type Storage struct {
	S3Client *s3.Client
	Bucket   string
}

type Resolver struct {
	URL string
}

func (r *Resolver) ResolveEndpoint(_ context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	endpointURL, err := url.Parse(r.URL)
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	endpointURL.Path += "/" + *params.Bucket
	return smithyendpoints.Endpoint{URI: *endpointURL}, nil
}

func NewStorage(url string, key string, secret string, bucket string) *Storage {
	client := s3.New(s3.Options{
		EndpointResolverV2: &Resolver{URL: url},
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     key,
				SecretAccessKey: secret,
			}, nil
		}),
	})

	return &Storage{
		S3Client: client,
		Bucket:   bucket,
	}
}

func strPtr(s string) *string {
	return &s
}

func (s Storage) UploadFile(ctx context.Context, objectKey string, file []byte) error {
	reader := bytes.NewReader(file)
	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: strPtr(s.Bucket),
		Key:    strPtr(objectKey),
		Body:   reader,
	})
	return err
}
