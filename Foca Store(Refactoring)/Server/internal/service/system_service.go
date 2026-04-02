package service

import (
	"voca-store/internal/helper"
	"voca-store/seeders"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SystemService interface {
	SeedRoles() error
	SeedAdmin() error
	SeedUsers() error
	SeedProducts() error
	SeedProductsFromAssets() error
	SyncAssetProducts() error
	SeedCoupons() error
	SeedAll() error
	SeedAllWithProducts() error
	ResetDatabase() error
	ResetDatabaseWithProducts() error
	ResetDatabasePreserveCatalog() error
	MigrateAll() error
	ResetRedis() error
	DeleteAllCloudinaryAssets() error
	GenerateSecret() (string, error)
}

type systemService struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewSystemService(db *gorm.DB, rdb *redis.Client) SystemService {
	return &systemService{db, rdb}
}

func (s *systemService) SeedRoles() error {
	return seeders.SeedRoles(s.db)
}

func (s *systemService) SeedAdmin() error {
	return seeders.SeedAdmin(s.db)
}

func (s *systemService) SeedUsers() error {
	return seeders.SeedUsers(s.db)
}

func (s *systemService) SeedProducts() error {
	return seeders.SeedProducts(s.db)
}

func (s *systemService) SeedProductsFromAssets() error {
	return seeders.SeedProductsFromAssets(s.db)
}

func (s *systemService) SyncAssetProducts() error {
	return seeders.SyncAssetProductsWithDefaultSeed(s.db)
}

func (s *systemService) SeedCoupons() error {
	return seeders.SeedCoupons(s.db)
}

func (s *systemService) SeedAll() error {
	if err := seeders.SeedRoles(s.db); err != nil { return err }
	if err := seeders.SeedAdmin(s.db); err != nil { return err }
	if err := seeders.SeedCategories(s.db); err != nil { return err }
	if err := seeders.SeedUsers(s.db); err != nil { return err }
	if err := seeders.SeedCoupons(s.db); err != nil { return err }
	return nil
}

func (s *systemService) SeedAllWithProducts() error {
	if err := s.SeedAll(); err != nil { return err }
	return seeders.SeedProducts(s.db)
}

func (s *systemService) ResetDatabase() error {
	return seeders.ResetDatabase(s.db, s.rdb)
}

func (s *systemService) ResetDatabaseWithProducts() error {
	return seeders.ResetDatabaseWithProduct(s.db, s.rdb)
}

func (s *systemService) ResetDatabasePreserveCatalog() error {
	return seeders.ResetDatabasePreserveProductsAndCategories(s.db, s.rdb)
}

func (s *systemService) MigrateAll() error {
	return seeders.MigrateAll(s.db)
}

func (s *systemService) ResetRedis() error {
	return seeders.ResetRedis(s.rdb)
}

func (s *systemService) DeleteAllCloudinaryAssets() error {
	return helper.DeleteAllAssets()
}

func (s *systemService) GenerateSecret() (string, error) {
	return helper.GenerateRandomSecret(32)
}
