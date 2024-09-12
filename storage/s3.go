package storage

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hantbk/vts-backup/helper"
	"github.com/hantbk/vts-backup/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 - Amazon S3 storage
//
// type: s3
// bucket: backup-test
// region: us-east-1
// path: backups
// access_key_id: your-access-key-id
// secret_access_key: your-secret-access-key
// max_retries: 5
// storage_class:
// timeout: 300
type S3 struct {
	Base
	Service      string
	bucket       string
	path         string
	client       *s3manager.Uploader
	storageClass string
	awsCfg       *aws.Config
}

func (s S3) providerName() string {

	switch s.Service {
	case "s3":
		return "AWS S3"
	case "minio":
		return "Minio"
	}
	return "AWS S3"
}

func (s S3) defaultRegion() string {
	switch s.Service {
	case "s3":
		return "us-east-1"
	case "minio":
		return "us-east-1"
	}
	return "us-east-1"
}

func (s S3) defaultEndpoint() *string {

	return aws.String("")
}

func (s *S3) defaultStorageClass() string {
	switch s.Service {
	case "s3":
		return "STANDARD_IA"
	case "minio":
		return ""
	}

	return ""
}

func (s *S3) forcePathStyle() bool {
	switch s.Service {
	case "tos", "oss":
		return false
	default:
		return true
	}
}

func (s *S3) init() {
	if len(s.Service) == 0 {
		s.Service = "s3"
	}

	s.viper.SetDefault("region", s.defaultRegion())
	s.viper.SetDefault("endpoint", s.defaultEndpoint())
	s.viper.SetDefault("max_retries", 3)
	s.viper.SetDefault("timeout", "300")
	s.viper.SetDefault("storage_class", s.defaultStorageClass())
}

func (s *S3) open() (err error) {
	s.init()

	cfg := aws.NewConfig()
	endpoint := s.viper.GetString("endpoint")

	if len(endpoint) > 0 {
		cfg.Endpoint = aws.String(endpoint)
	}

	cfg.S3ForcePathStyle = aws.Bool(s.forcePathStyle())
	if s.viper.IsSet("force_path_style") {
		cfg.S3ForcePathStyle = aws.Bool(s.viper.GetBool("force_path_style"))
	}

	accessKeyId := s.viper.GetString("access_key_id")
	secretAccessKey := s.viper.GetString("secret_access_key")
	if len(secretAccessKey) == 0 {
		secretAccessKey = s.viper.GetString("access_key_secret")
	}

	cfg.Credentials = credentials.NewStaticCredentials(
		accessKeyId,
		secretAccessKey,
		s.viper.GetString("token"),
	)

	cfg.Region = aws.String(s.viper.GetString("region"))
	cfg.MaxRetries = aws.Int(s.viper.GetInt("max_retries"))

	s.bucket = s.viper.GetString("bucket")
	s.path = s.viper.GetString("path")

	// Only present storage_class when it is set.
	if len(s.viper.GetString("storage_class")) > 0 {
		s.storageClass = s.viper.GetString("storage_class")
	}

	timeout := s.viper.GetInt("timeout")
	uploadTimeoutDuration := time.Duration(timeout) * time.Second

	httpClient := &http.Client{Timeout: uploadTimeoutDuration}
	cfg.HTTPClient = httpClient
	s.awsCfg = cfg

	sess := session.Must(session.NewSession(s.awsCfg))
	s.client = s3manager.NewUploader(sess)

	return
}

func (s *S3) close() {
}

func (s *S3) upload(fileKey string) (err error) {
	logger := logger.Tag(s.providerName())

	var fileKeys []string
	if len(s.fileKeys) != 0 {
		// directory
		// 2022.12.04.07.09.47/2022.12.04.07.09.47.tar.xz-000
		fileKeys = s.fileKeys
	} else {
		// file
		// 2022.12.04.07.09.25.tar.xz
		fileKeys = append(fileKeys, fileKey)
	}

	for _, key := range fileKeys {
		sourcePath := filepath.Join(filepath.Dir(s.archivePath), key)
		remotePath := filepath.Join(s.path, key)

		f, err := os.Open(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to open file %q, %v", sourcePath, err)
		}
		defer f.Close()

		progress := helper.NewProgressBar(logger, f)

		input := &s3manager.UploadInput{
			Bucket:       aws.String(s.bucket),
			Key:          aws.String(remotePath),
			Body:         progress.Reader,
			StorageClass: aws.String(s.storageClass),
		}

		result, err := s.client.Upload(input, func(uploader *s3manager.Uploader) {
			// set the part size as low as possible to avoid timeouts and aborts
			// also set concurrency to 1 for the same reason
			var partSize int64 = 64 * 1024 * 1024 // 64MiB
			maxParts := math.Ceil(float64(progress.FileLength / partSize))

			// 10000 parts is the limit for AWS S3. If the resulting number of parts would exceed that limit, increase the
			// part size as much as needed but as little possible
			if maxParts > 10000 {
				partSize = int64(math.Ceil(float64(progress.FileLength) / 10000))
			}

			uploader.Concurrency = 1
			uploader.LeavePartsOnError = false
			uploader.PartSize = partSize
		})

		if err != nil {
			return progress.Errorf("%v", err)
		}

		progress.Done(result.Location)

		if s.Service == "s3" {
			logger.Info("=>", fmt.Sprintf("s3://%s/%s", s.bucket, remotePath))
		}
	}

	return nil
}

func (s *S3) delete(fileKey string) (err error) {
	remotePath := filepath.Join(s.path, fileKey)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(remotePath),
	}
	_, err = s.client.S3.DeleteObject(input)
	return
}

// List the objects in the bucket with the prefix = parent
func (s *S3) list(parent string) ([]FileItem, error) {
	remotePath := filepath.Join(s.path, parent)
	continueToken := ""
	var items []FileItem

	for {
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.bucket),
			Prefix:            aws.String(remotePath),
			ContinuationToken: aws.String(continueToken),
		}

		result, err := s.client.S3.ListObjectsV2(input)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects, %v", err)
		}

		for _, object := range result.Contents {
			items = append(items, FileItem{
				Filename:     *object.Key,
				Size:         *object.Size,
				LastModified: *object.LastModified,
			})
		}

		if *result.IsTruncated {
			continueToken = *result.NextContinuationToken
		} else {
			break
		}
	}

	return items, nil
}

// Get the object download URL by fileKey (include remote_path)
func (s *S3) download(fileKey string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileKey),
	}

	req, _ := s.client.S3.GetObjectRequest(input)
	url, err := req.Presign(1 * time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to sign request, %v", err)
	}

	return url, nil
}
