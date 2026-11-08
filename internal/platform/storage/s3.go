package storage

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// S3Config holds everything needed to reach one bucket.
type S3Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	PathStyle bool
}

// S3 is an ObjectStore backed by an S3-compatible API (Garage, AWS S3, MinIO).
type S3 struct {
	client *s3.Client
	bucket string
}

// NewS3 builds a client with static credentials. Path-style addressing and a
// custom endpoint are what Garage and most self-hosted stores expect.
func NewS3(ctx context.Context, cfg S3Config) (*S3, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey, cfg.SecretKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("storage: load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = cfg.PathStyle
		// Garage does not implement every checksum header the SDK adds by
		// default, so only send one when the operation requires it.
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		// Reverse proxies (Cloudflare in front of Garage, for example) rewrite
		// Accept-Encoding at the edge. The SDK signs that header by default, so
		// a rewritten value turns every SigV4 signature invalid.
		o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
			return stack.Finalize.Insert(&dropAcceptEncoding{}, "DisableAcceptEncodingGzip", middleware.After)
		})
	})

	return &S3{client: client, bucket: cfg.Bucket}, nil
}

// dropAcceptEncoding removes Accept-Encoding before the request is signed. It
// runs right after the SDK's gzip-handling middleware sets the header.
type dropAcceptEncoding struct{}

func (*dropAcceptEncoding) ID() string { return "DropAcceptEncoding" }

func (*dropAcceptEncoding) HandleFinalize(
	ctx context.Context,
	in middleware.FinalizeInput,
	next middleware.FinalizeHandler,
) (middleware.FinalizeOutput, middleware.Metadata, error) {
	if req, ok := in.Request.(*smithyhttp.Request); ok {
		req.Header.Del("Accept-Encoding")
	}
	return next.HandleFinalize(ctx, in)
}

// Save uploads the body and returns the number of bytes read from it.
func (s *S3) Save(ctx context.Context, in SaveInput) (string, int64, error) {
	counter := &countingReader{reader: in.Body}
	contentType := in.ContentType
	if contentType == "" {
		contentType = DefaultContentType
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(in.Path),
		Body:        counter,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", 0, fmt.Errorf("storage: put object %q: %w", in.Path, err)
	}
	return in.Path, counter.size, nil
}

// Read returns the object body. The caller owns closing it.
func (s *S3) Read(ctx context.Context, objectPath string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectPath),
	})
	if err != nil {
		if isNotFound(err) {
			return nil, ErrObjectNotFound
		}
		return nil, fmt.Errorf("storage: get object %q: %w", objectPath, err)
	}
	return out.Body, nil
}

// Delete removes the object. Deleting a missing object is not an error, which
// matches S3 semantics.
func (s *S3) Delete(ctx context.Context, objectPath string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(objectPath),
	})
	if err != nil {
		return fmt.Errorf("storage: delete object %q: %w", objectPath, err)
	}
	return nil
}

func isNotFound(err error) bool {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	switch apiErr.ErrorCode() {
	case "NoSuchKey", "NotFound":
		return true
	}
	return false
}

type countingReader struct {
	reader io.Reader
	size   int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	c.size += int64(n)
	return n, err
}
