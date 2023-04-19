package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := getEnvDefault("DATABASE_URL", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		getEnvDefault("DB_HOST", "localhost"),
		getEnvDefault("DB_USER", "postgres"),
		getEnvDefault("DB_PASSWORD", "postgres"),
		getEnvDefault("DB_NAME", "postgres"),
		getEnvDefault("DB_PORT", "5432")))

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = database.AutoMigrate(
		&MultiChoiceQuestion{},
		&MultiChoiceAnswer{},
		&ShortAnswerQuestion{},
		&FormativeAssessment{},
		&Tag{})
	if err != nil {
		return
	}
	DB = database
}
func getEnvDefault(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
