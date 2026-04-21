package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"server/internal/helper"
)

type MinIOStorage struct {
	client   *minio.Client
	bucket   string
	baseURL  string
	endpoint string
	useSSL   bool
}

func NewMinIO(cfg helper.Config) (*MinIOStorage, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.MinIOBucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &MinIOStorage{
		client:   client,
		bucket:   cfg.MinIOBucket,
		baseURL:  cfg.StorageBaseURL,
		endpoint: cfg.MinIOEndpoint,
		useSSL:   cfg.MinIOUseSSL,
	}, nil
}

func (s *MinIOStorage) UploadDirectory(prefix string, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		objectName := filepath.ToSlash(filepath.Join(prefix, relative))
		contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(path)))
		_, err = s.client.FPutObject(context.Background(), s.bucket, objectName, path, minio.PutObjectOptions{
			ContentType: contentType,
		})
		return err
	})
}

func (s *MinIOStorage) BuildPrefixURL(prefix string) string {
	if s.baseURL != "" {
		return fmt.Sprintf("%s/%s", s.baseURL, prefix)
	}
	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.endpoint, s.bucket, prefix)
}

func (s *MinIOStorage) PresignedObjectURL(objectName string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(context.Background(), s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *MinIOStorage) ExtractObjectPrefix(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	if !strings.Contains(trimmed, "://") {
		return strings.Trim(trimmed, "/")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return strings.Trim(trimmed, "/")
	}

	path := strings.Trim(parsed.Path, "/")
	if s.bucket != "" && strings.HasPrefix(path, s.bucket+"/") {
		return strings.TrimPrefix(path, s.bucket+"/")
	}
	return path
}

func (s *MinIOStorage) GetObject(objectName string) (*minio.Object, minio.ObjectInfo, error) {
	object, err := s.client.GetObject(context.Background(), s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}

	info, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, minio.ObjectInfo{}, err
	}

	return object, info, nil
}

func (s *MinIOStorage) CloseObject(object io.Closer) error {
	if object == nil {
		return nil
	}
	return object.Close()
}

func (s *MinIOStorage) ClearBucket() error {
	ctx := context.Background()
	objects := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for removeErr := range s.client.RemoveObjects(ctx, s.bucket, objects, minio.RemoveObjectsOptions{}) {
		if removeErr.Err != nil {
			return removeErr.Err
		}
	}

	return nil
}
