package router

import (
	"github.com/cemayan/event-scraper/api/src/database"
	"github.com/cemayan/event-scraper/api/src/middleware"
	"github.com/cemayan/event-scraper/api/src/repo"
	"github.com/cemayan/event-scraper/api/src/service"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
	"time"
)

var cb *gobreaker.CircuitBreaker

// SetupRoutes creates the fiber's routes
// api/v1 is root group.
// Before the reach services interface is configured
func SetupRoutes(app *fiber.App, log *log.Logger) {

	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	//Resty http client
	httpClient := resty.New()

	// CircuitBreaker
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.Timeout = 10 * time.Second
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb = gobreaker.NewCircuitBreaker(st)

	var _middleware middleware.Middleware
	_middleware = middleware.NewMiddleware(httpClient, cb, log)

	var eventRepo repo.EventRepository
	eventRepo = repo.NewEventRepo(database.GetDB(), log)

	var eventSvc service.EventService
	eventSvc = service.NewEventService(eventRepo, log)

	v1.Get("/health", eventSvc.HealthCheck)

	user := v1.Group("/event")
	user.Get("/provider/:provider", _middleware.Protected(), _middleware.ValidUserId, eventSvc.GetByProvider)

}
