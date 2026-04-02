package helper

import (
	"fmt"
	"regexp"
	"strings"
	"voca-store/internal/domain/models"

	"gorm.io/gorm"
)

func GenerateSlug(input string) string {
	slug := strings.ToLower(input)
	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	slug = reg.ReplaceAllString(slug, "")

	regDash := regexp.MustCompile(`-+`)
	slug = regDash.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}

func GenerateUniqueSlug(db *gorm.DB, name string) (string, error) {
	baseSlug := GenerateSlug(name)
	slug := baseSlug
	counter := 1

	for {
		var count int64
		err := db.Model(&models.Product{}).
			Where("slug = ?", slug).
			Count(&count).Error
		
		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug, nil
}

func GenerateUniqueCategorySlug(db *gorm.DB, name string) (string, error) {
	baseSlug := GenerateSlug(name)
	slug := baseSlug
	counter := 1

	for {
		var count int64

		err := db.Model(&models.Category{}).Where("slug = ?", slug).Count(&count).Error

		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}

		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	return slug, nil
}
