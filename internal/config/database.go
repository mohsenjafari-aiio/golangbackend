package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	orderDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/order/domain"
	productDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/product/domain"
	userDomain "github.com/mohsenjafari-aiio/aiiobackend/internal/user/domain"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "aiio_backend"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func (config *DatabaseConfig) BuildDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
}

func ConnectDatabase() (*gorm.DB, error) {
	config := GetDatabaseConfig()
	dsn := config.BuildDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schemas
	err = db.AutoMigrate(
		&userDomain.User{},
		&productDomain.Product{},
		&orderDomain.Order{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
