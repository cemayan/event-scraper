package main

import (
	"github.com/cemayan/event-scraper/config/user"
	"github.com/cemayan/event-scraper/internal/user/database"
	"github.com/cemayan/event-scraper/internal/user/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var _log *logrus.Logger
var app = fiber.New()
var configs *user.AppConfig
var v *viper.Viper
var dbHandler database.DBHandler

func init() {
	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout

	v = viper.New()
	_configs := user.NewConfig(v)

	env := os.Getenv("ENV")
	appConfig, err := _configs.GetConfig(env)
	configs = appConfig
	if err != nil {
		return
	}

	//Postresql connection
	dbHandler = database.NewDbHandler(configs)
	dbHandler.ConnectDB()

}

func main() {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	_log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	router.SetupRoutes(app, _log, configs)

	err := app.Listen(":8089")
	if err != nil {
		return
	}

}
