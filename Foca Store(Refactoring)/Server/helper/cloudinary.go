package helper

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cloudinaryInstance *cloudinary.Cloudinary

// Struct untuk menampung hasil upload agar bisa diakses di handler
type UploadResult struct {
	SecureURL string
	PublicID  string
}

// Helper untuk menangani *bool
func ptrBool(b bool) *bool {
	return &b
}

func InitCloudinary() error {
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL != "" {
		var err error
		cloudinaryInstance, err = cloudinary.NewFromURL(cloudinaryURL)
		if err != nil {
			return fmt.Errorf("failed to initialize cloudinary from URL: %w", err)
		}
		return nil
	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return errors.New("cloudinary credentials not found")
	}

	var err error
	cloudinaryInstance, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("failed to initialize cloudinary from params: %w", err)
	}

	return nil
}

func UploadFile(file interface{}, folder string) (*UploadResult, error) {
	if cloudinaryInstance == nil {
		return nil, errors.New("cloudinary not initialized")
	}

	uploadParam := uploader.UploadParams{
		Folder: folder,
		UseFilename: ptrBool(true),
		UniqueFilename: ptrBool(true),
		Overwrite: ptrBool(false),
		ResourceType: "image",
	}

	resp, err := cloudinaryInstance.Upload.Upload(context.Background(), file, uploadParam)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		SecureURL: resp.SecureURL,
		PublicID: resp.PublicID,
	}, nil
}

func DeleteImage(publicID string) error {

	if cloudinaryInstance == nil {
		return errors.New("cloudinary not initialized")
	}

	resp, err := cloudinaryInstance.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID:   publicID,
			Invalidate: ptrBool(true),
			ResourceType: "image",
			Type: "upload",
		},
	)

	if err != nil {
		return fmt.Errorf("cloudinary delete error: %w", err)
	}

	// VALIDASI WAJIB
	if resp.Result != "ok" {
		return fmt.Errorf("cloudinary deletion failed: %s", resp.Result)
	}

	fmt.Println("Cloudinary deleted:", publicID)

	return nil
}

func DeleteAllAssets() error {
	if cloudinaryInstance == nil {
		return errors.New("cloudinary not initialized")
	}

	ctx := context.Background()
	nextCursor := ""

	for {
		result, err := cloudinaryInstance.Admin.Assets(ctx, admin.AssetsParams{
			MaxResults: 500,
			NextCursor: nextCursor,
		})

		if err != nil {
			return fmt.Errorf("failed to fetch assets: %w", err)
		}

		if len(result.Assets) == 0 {
			break
		}

		var publicIDs []string
		for _, asset := range result.Assets {
			publicIDs = append(publicIDs, asset.PublicID)
		}

		_, err = cloudinaryInstance.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{
			PublicIDs: publicIDs,
			Invalidate: ptrBool(true),
		})

		if err != nil {
			return fmt.Errorf("failed to delete assets: %w", err)
		}

		nextCursor = result.NextCursor
	}

	return nil
}