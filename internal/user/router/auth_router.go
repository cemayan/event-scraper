package router

import (
	"github.com/cemayan/event-scraper/config/user"
	"github.com/cemayan/event-scraper/internal/user/database"
	"github.com/cemayan/event-scraper/internal/user/repo"
	"github.com/cemayan/event-scraper/internal/user/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

// SetupAuthRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupAuthRoutes(app *fiber.App, log *log.Logger, configs *user.AppConfig) {

	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	userRepo := repo.NewUserRepo(database.DB, log)

	var authSvc service.AuthService
	authSvc = service.NewAuthService(userRepo, log, configs)

	v1.Get("/health", authSvc.HealthCheck)

	//Auth
	auth := v1.Group("/auth")
	auth.Post("/getToken", authSvc.Login)
	auth.Post("/validateToken", authSvc.ValidToken)

}
