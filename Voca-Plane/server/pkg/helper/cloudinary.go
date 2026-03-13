package helper

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(file *multipart.FileHeader) (string, string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	result, err := cld.Upload.Upload(context.Background(), src, uploader.UploadParams{
		Folder: "voca-plane/airlines",
	})
	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}

func DeleteImage(publicID string) error {

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return err
	}

	_, err = cld.Upload.Destroy(
		context.Background(),
		uploader.DestroyParams{
			PublicID: publicID,
		},
	)

	return err
}