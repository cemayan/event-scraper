package service

import "github.com/gofiber/fiber/v2"

type EventService interface {
	GetByProvider(c *fiber.Ctx) error
	DeleteByProvider(provider string)
	HealthCheck(c *fiber.Ctx) error
}
