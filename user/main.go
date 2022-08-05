package main

import (
	"github.com/cemayan/event-scraper/user/src/config"
	"github.com/cemayan/event-scraper/user/src/database"
	"github.com/cemayan/event-scraper/user/src/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"os"
)

var _log *logrus.Logger
var app = fiber.New()

func init() {
	//logrus init
	_log = logrus.New()
	_log.Out = os.Stdout
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

	configs := config.GetConfig()
	database.ConnectDB(configs)

	router.SetupRoutes(app, _log)

	err := app.Listen(":8089")
	if err != nil {
		return
	}

}
