package router

import (
	"github.com/cemayan/event-scraper/config/user"
	"github.com/cemayan/event-scraper/internal/user/database"
	"github.com/cemayan/event-scraper/internal/user/middleware"
	"github.com/cemayan/event-scraper/internal/user/repo"
	service2 "github.com/cemayan/event-scraper/internal/user/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

// SetupRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupRoutes(app *fiber.App, log *log.Logger, configs *user.AppConfig) {

	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	userRepo := repo.NewUserRepo(database.DB, log)

	var authSvc service2.AuthService
	authSvc = service2.NewAuthService(userRepo, log, configs)

	var userSvc service2.UserService
	userSvc = service2.NewUserService(userRepo, authSvc, log, configs)

	v1.Get("/health", userSvc.HealthCheck)

	user := v1.Group("/user")
	user.Get("/:id", middleware.Protected(configs), userSvc.GetUser)
	user.Post("/", userSvc.CreateUser)
	user.Put("/:id", middleware.Protected(configs), userSvc.UpdateUser)
	user.Delete("/:id", middleware.Protected(configs), userSvc.DeleteUser)

}
