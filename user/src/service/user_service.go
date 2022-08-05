package service

import (
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	hashPassword(password string) (string, error)
	GetUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	HealthCheck(c *fiber.Ctx) error
}
