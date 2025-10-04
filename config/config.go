package config

import (
	"os"

	"github.com/MujiRahman/golang-simple-note/internal/helper"
	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string `yaml:"db_user"`
	DBPassword string `yaml:"db_password"`
	DBHost     string `yaml:"db_host"`
	DBPort     string `yaml:"db_port"`
	DBName     string `yaml:"db_name"`
	JWTSecret  string
	TokenTTL   int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	helper.LogIfError(err, "Warning: .env file not found, using system env")

	return &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		TokenTTL:   3600, // Default token TTL in seconds
	}
}
