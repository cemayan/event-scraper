package middleware

import (
	"github.com/cemayan/event-scraper/user/src/config"
	"github.com/cemayan/event-scraper/user/src/model"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	configs := config.GetConfig()
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(configs.SECRET),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(400).JSON(model.Response{
			Message: "Missing or malformed JWT",
		})

	}
	return c.Status(401).JSON(model.Response{
		Message: "Invalid or expired JWT",
	})
}
