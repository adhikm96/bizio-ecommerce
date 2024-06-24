package database

import (
	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"strconv"
)

var db *gorm.DB

func NewDatabaseConnection() *gorm.DB {
	if db != nil {
		return db
	}

	dbUrl := "postgresql://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	sqlDB, _ := db.DB()

	maxOpenConn, _ := strconv.Atoi(os.Getenv("DB_POOL_MAX_SIZE"))
	sqlDB.SetMaxOpenConns(maxOpenConn)

	if err != nil {
		slog.Error(err.Error())
	}

	return db
}

func MigrateDBSchema() {
	db = NewDatabaseConnection()
	err := db.AutoMigrate(&model.User{})

	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info("DB migrated")
}
