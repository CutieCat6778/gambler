package database

import (
	"fmt"
	"gambler/backend/database/models"
	"gambler/backend/tools"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "[DATABASE]\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second * 3,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	Database, err := gorm.Open(postgres.Open(tools.DATABASE), &gorm.Config{
		TranslateError: true,
		Logger:         newLogger,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("[DATABASE] Database connected")

	Database.AutoMigrate(&models.User{})

	return Database
}
