package database

import (
	"fmt"
	"log"
	"marketplace/config"
	"marketplace/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedRoles()
	log.Println("Database connected and migrated successfully")
}

func seedRoles() {
	roles := []models.Role{
		{Name: "admin", Description: "Administrator"},
		{Name: "seller", Description: "Seller"},
		{Name: "buyer", Description: "Buyer"},
	}
	for _, role := range roles {
		var existing models.Role
		DB.Where("name = ?", role.Name).First(&existing)
		if existing.ID == 0 {
			DB.Create(&role)
		}
	}
}
