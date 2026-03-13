package config

import (
	"github.com/cloudinary/cloudinary-go"
)

func NewClaudinaryClient(cfg *CloudinaryConfig) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	if err != nil {
		return nil, err
	}

	return cld, err
}