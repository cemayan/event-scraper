package database

import (
	"context"
	"fmt"
	"github.com/cemayan/event-scraper-common/protos"
	"github.com/cemayan/event-scraper/api/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// ConnectDB  serves to connect to db
// When DB connection is successful then model migration is started
func ConnectDB(config config.Config) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s  user=%s password=%s  dbname=%s sslmode=disable ",
			config.DB_HOST,
			config.DB_PORT,
			config.DB_USER,
			config.DB_PASSWORD,
			config.DB_NAME),
	}), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	var event protos.Event
	orm, err := event.ToORM(context.Background())
	if err != nil {
		return
	}

	if os.Getenv("ENV") == "test" {
		// ConnectDBForTesting  serves to connect to db for Testing
		// When DB connection is successful then model migration is started
		db.Migrator().DropTable(protos.Event{})
		db.AutoMigrate(orm)
		fmt.Println("Database Migrated")
	} else {
		db.AutoMigrate(orm)
		fmt.Println("Database Migrated")
	}

	DB = db
}

// GetDB gives gorm DB connector
// If connection is nil then is created new connection
func GetDB() *gorm.DB {
	if DB == nil {
		configs := config.GetConfig()
		ConnectDB(configs)
	}
	return DB
}
