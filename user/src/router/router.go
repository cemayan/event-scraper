package router

import (
	"github.com/cemayan/event-scraper/user/src/database"
	"github.com/cemayan/event-scraper/user/src/middleware"
	"github.com/cemayan/event-scraper/user/src/repo"
	"github.com/cemayan/event-scraper/user/src/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

// SetupRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupRoutes(app *fiber.App, log *log.Logger) {

	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	userRepo := repo.NewUserRepo(database.DB, log)

	var authSvc service.AuthService
	authSvc = service.NewAuthService(userRepo, log)

	var userSvc service.UserService
	userSvc = service.NewUserService(userRepo, authSvc, log)

	v1.Get("/health", userSvc.HealthCheck)

	user := v1.Group("/user")
	user.Get("/:id", middleware.Protected(), userSvc.GetUser)
	user.Post("/", userSvc.CreateUser)
	user.Put("/:id", middleware.Protected(), userSvc.UpdateUser)
	user.Delete("/:id", middleware.Protected(), userSvc.DeleteUser)

}
