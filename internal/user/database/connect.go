package database

import (
	"fmt"
	"github.com/cemayan/event-scraper/config/user"
	"github.com/cemayan/event-scraper/internal/user/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DBHandler interface {
	ConnectDB()
}

type DBService struct {
	configs *user.AppConfig
}

// ConnectDB  serves to connect to db
// When DB connection is successful then model migration is started
func (d DBService) ConnectDB() {

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
			d.configs.Postgresql.HOST,
			d.configs.Postgresql.PORT,
			d.configs.Postgresql.USER,
			d.configs.Postgresql.PASSWORD,
			d.configs.Postgresql.NAME),
	}), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	if os.Getenv("ENV") == "test" {
		// ConnectDBForTesting  serves to connect to db for Testing
		// When DB connection is successful then model migration is started
		db.Migrator().DropTable(&model.User{})
		db.AutoMigrate(&model.User{})
		fmt.Println("Database Migrated")
	} else {
		db.AutoMigrate(&model.User{})
		fmt.Println("Database Migrated")
	}

	DB = db
}

func NewDbHandler(configs *user.AppConfig) DBHandler {
	return &DBService{configs: configs}
}
