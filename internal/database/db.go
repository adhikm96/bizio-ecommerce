package database

import (
	"gorm.io/gorm/logger"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/Digital-AIR/bizio-ecommerce/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewDatabaseConnection() *gorm.DB {
	if db != nil {
		return db
	}

	dbUrl := "postgresql://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,        // Don't include params in the SQL log
				Colorful:                  false,       // Disable color
			},
		),
	})

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
	err := db.AutoMigrate(
		&model.User{},
		&model.Notification{},
		&model.Attribute{},
		&model.AttributeValue{},
		&model.Discount{},
		&model.Category{},
		&model.Brand{},
		&model.Product{},
		&model.ProductVariant{},
		&model.Inventory{},
		&model.Address{},
		&model.Cart{},
		&model.Promotion{},
		&model.Order{},
		&model.Payment{},
		&model.CartItem{},
		&model.VariantAttribute{},
		&model.OrderItem{},
		&model.ProductPromotion{},
		&model.Review{},
	)

	if err != nil {
		slog.Error(err.Error())
		panic(err.Error())
	}

	slog.Info("DB migrated")
}
