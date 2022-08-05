package service

import (
	"github.com/cemayan/event-scraper/api/src/repo"
	"github.com/cemayan/event-scraper/api/src/utils"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// A EventSvc  contains the required dependencies for this service
type EventSvc struct {
	repository repo.EventRepository
	log        *log.Logger
}

func (e EventSvc) HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).Send([]byte("OK!"))
}

// GetByProvider returns Events based on given provider
func (e EventSvc) GetByProvider(c *fiber.Ctx) error {

	provider, _ := c.ParamsInt("provider")
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	tt := utils.Provider(provider)
	return c.JSON(e.repository.GetByProvider(tt, page, pageSize))
}

// DeleteByProvider removes based on given provider
func (e EventSvc) DeleteByProvider(provider string) {
	e.repository.DeleteByProvider(provider)
}

func NewEventService(repository repo.EventRepository, log *log.Logger) EventService {
	return &EventSvc{repository: repository, log: log}
}
