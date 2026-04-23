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
	client          *minio.Client
	bucket          string
	thumbnailBucket string
	profileBucket   string
	baseURL         string
	endpoint        string
	useSSL          bool
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
	for _, bucket := range []string{cfg.MinIOBucket, cfg.MinIOThumbnailBucket, cfg.MinIOProfileBucket} {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return nil, err
		}
		if !exists {
			if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return nil, err
			}
		}
	}

	return &MinIOStorage{
		client:          client,
		bucket:          cfg.MinIOBucket,
		thumbnailBucket: cfg.MinIOThumbnailBucket,
		profileBucket:   cfg.MinIOProfileBucket,
		baseURL:         cfg.StorageBaseURL,
		endpoint:        cfg.MinIOEndpoint,
		useSSL:          cfg.MinIOUseSSL,
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

func (s *MinIOStorage) UploadFile(bucket string, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(context.Background(), bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *MinIOStorage) BuildPrefixURL(prefix string) string {
	return s.buildBucketObjectURL(s.bucket, prefix)
}

func (s *MinIOStorage) BuildThumbnailURL(objectName string) string {
	return s.buildBucketObjectURL(s.thumbnailBucket, objectName)
}

func (s *MinIOStorage) BuildProfileURL(objectName string) string {
	return s.buildBucketObjectURL(s.profileBucket, objectName)
}

func (s *MinIOStorage) ThumbnailBucket() string {
	return s.thumbnailBucket
}

func (s *MinIOStorage) ProfileBucket() string {
	return s.profileBucket
}

func (s *MinIOStorage) buildBucketObjectURL(bucket string, objectName string) string {
	if s.baseURL != "" {
		return fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.baseURL, "/"), bucket, strings.TrimLeft(objectName, "/"))
	}
	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.endpoint, bucket, strings.TrimLeft(objectName, "/"))
}

func (s *MinIOStorage) PresignedObjectURL(objectName string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(context.Background(), s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *MinIOStorage) ExtractObjectPrefix(value string) string {
	return s.extractObjectName(s.bucket, value)
}

func (s *MinIOStorage) ExtractProfileObjectName(value string) string {
	return s.extractObjectName(s.profileBucket, value)
}

func (s *MinIOStorage) extractObjectName(bucket string, value string) string {
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
	if bucket != "" && strings.HasPrefix(path, bucket+"/") {
		return strings.TrimPrefix(path, bucket+"/")
	}
	return path
}

func (s *MinIOStorage) GetObject(objectName string) (*minio.Object, minio.ObjectInfo, error) {
	return s.getObjectFromBucket(s.bucket, objectName)
}

func (s *MinIOStorage) GetThumbnailObject(objectName string) (*minio.Object, minio.ObjectInfo, error) {
	return s.getObjectFromBucket(s.thumbnailBucket, objectName)
}

func (s *MinIOStorage) GetProfileObject(objectName string) (*minio.Object, minio.ObjectInfo, error) {
	return s.getObjectFromBucket(s.profileBucket, objectName)
}

func (s *MinIOStorage) getObjectFromBucket(bucket string, objectName string) (*minio.Object, minio.ObjectInfo, error) {
	object, err := s.client.GetObject(context.Background(), bucket, objectName, minio.GetObjectOptions{})
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
	return s.clearBucket(s.bucket)
}

func (s *MinIOStorage) ClearThumbnailBucket() error {
	return s.clearBucket(s.thumbnailBucket)
}

func (s *MinIOStorage) ClearProfileBucket() error {
	return s.clearBucket(s.profileBucket)
}

func (s *MinIOStorage) RemoveProfileObject(objectName string) error {
	return s.removeObject(s.profileBucket, objectName)
}

func (s *MinIOStorage) clearBucket(bucket string) error {
	ctx := context.Background()
	objects := s.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for removeErr := range s.client.RemoveObjects(ctx, bucket, objects, minio.RemoveObjectsOptions{}) {
		if removeErr.Err != nil {
			return removeErr.Err
		}
	}

	return nil
}

func (s *MinIOStorage) removeObject(bucket string, objectName string) error {
	if strings.TrimSpace(objectName) == "" {
		return nil
	}

	return s.client.RemoveObject(context.Background(), bucket, objectName, minio.RemoveObjectOptions{})
}
