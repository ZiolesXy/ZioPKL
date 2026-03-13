package config

import (
	"fmt"
	"log"
	"voca-plane/internal/domain/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func PostgresDatabase(cfg *config) *gorm.DB {
	fmt.Println("📊  Connecting Databases")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database; %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Airline{},
		&models.Airport{},
		&models.Flight{},
		&models.FlightClass{},
		&models.Seat{},
		&models.FlightSeat{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.PromoCode{},
	)

	log.Println("📊  Databases Succesfully Connected")
	return db
}

func MySQLDatabase(cfg *config) *gorm.DB {
	log.Println("📊 Connecting Databases")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Airline{},
		&models.Airport{},
		&models.Flight{},
		&models.FlightClass{},
		&models.Seat{},
		&models.FlightSeat{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.PromoCode{},
	)

	log.Println("📊 Databases Succesfully Connected")
	return db
}