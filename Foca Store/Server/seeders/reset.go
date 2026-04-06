package seeders

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ctx = context.Background()

func ResetRedis(rdb *redis.Client) error {
	return rdb.FlushAll(ctx).Err()
}

func ResetDatabase(db *gorm.DB, rdb *redis.Client) error {
	if err := ResetRedis(rdb); err != nil {
		return err
	}
	if err := DropAllTable(db); err != nil {
		return err
	}

	if err := MigrateAll(db); err != nil {
		return err
	}

	if err := SeedRoles(db); err != nil {
			return err
		}

	if err := SeedAdmin(db); err != nil {
		return err
	}

	if err := SeedCategories(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedCoupons(db); err != nil {
		return err
	}
	
	return nil
}

func ResetDatabaseWithProduct(db *gorm.DB, rdb *redis.Client) error {
	if err := ResetRedis(rdb); err != nil {
		return err
	}

	if err := DropAllTable(db); err != nil {
		return err
	}

	if err := MigrateAll(db); err != nil {
		return err
	}

	if err := SeedRoles(db); err != nil {
			return err
		}

	if err := SeedAdmin(db); err != nil {
		return err
	}

	if err := SeedCategories(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedCoupons(db); err != nil {
		return err
	}

	if err := SeedProducts(db); err != nil {
		return err
	}
	
	return nil
}

func ResetDatabasePreserveProductsAndCategories(db *gorm.DB, rdb *redis.Client) error {
	if err := ResetRedis(rdb); err != nil {
		return err
	}

	if err := DropTableExceptProductsAndCategories(db); err != nil {
		return err
	}

	if err := MigrateAll(db); err != nil {
		return err
	}

	if err := SeedRoles(db); err != nil {
		return err
	}

	if err := SeedAdmin(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedCoupons(db); err != nil {
		return err
	}

	if err := SyncAssetProductsWithDefaultSeed(db); err != nil {
		return err
	}

	return nil
}
