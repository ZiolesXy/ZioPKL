package repository

import (
	"voca-store/internal/domain/models"
	"voca-store/internal/domain/repository"
	"voca-store/internal/helper"

	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.AddressRepository {
	return &addressRepository{db}
}

func (r *addressRepository) Create(address *models.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) FindByID(id uint, userID uint) (models.Address, error) {
	var address models.Address
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&address).Error
	return address, err
}

func (r *addressRepository) FindByUID(uid string, userID uint) (models.Address, error) {
	var address models.Address
	err := r.db.Where("uid = ? AND user_id = ?", uid, userID).First(&address).Error
	return address, err
}

func (r *addressRepository) FindAll(userID uint) ([]models.Address, error) {
	var addresses []models.Address
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) Update(address *models.Address, updates map[string]interface{}) error {
	return r.db.Model(address).Updates(updates).Error
}

func (r *addressRepository) Delete(address *models.Address) error {
	return r.db.Delete(address).Error
}

func (r *addressRepository) UnsetPrimary(userID uint) error {
	return r.db.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_primary", false).Error
}

func (r *addressRepository) WithTransaction(fn func(repo repository.AddressRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewAddressRepository(tx))
	})
}

func (r *addressRepository) GenerateUID() (string, error) {
	return helper.NewGenerateAddressUID(r.db)
}
