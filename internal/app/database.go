package app

import (
	"fmt"
	"log"

	"github.com/MujiRahman/golang-simple-note/config"
	"github.com/MujiRahman/golang-simple-note/internal/helper"
	"github.com/MujiRahman/golang-simple-note/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connect struct {
	DB *gorm.DB
}

func NewDB(cfg *config.Config) *Connect {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	helper.LogFatalIfError(err, "failed to connect to MariaDB: %v")

	helper.Print("koneksi data base berhasil yey")
	// Auto migrate
	err = db.AutoMigrate(&model.Note{}, &model.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	return &Connect{DB: db}
}
