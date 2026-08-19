package db

import(
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(databaseURL, env string) (*gorm.DB, error) {

	logLevel := logger.Silent
	if env == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(databaseURL),  &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	//Ping() - bcz gorm.Open can succeed even if the DB is unreachable, depending on driver
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("database connection established")
	return db, nil
}