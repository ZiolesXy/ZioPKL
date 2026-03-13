package repository

import (
	"context"
	"vantaze/pkg/helper"
)

type UploadResult struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
	Folder   string `json:"folder"`
	Format   string `json:"format"`
	Bytes    int    `json:"bytes"`
}

type StorageRepository interface {
	Upload(ctx context.Context, file helper.FileData, folder string) (*UploadResult, error)
	Get(ctx context.Context, publicID string) (*UploadResult, error)
	ListByFolder(ctx context.Context, folder string) ([]UploadResult, error)
	Delete(ctx context.Context, publicID string) error
	DeleteByFolder(ctx context.Context, folder string) error
}