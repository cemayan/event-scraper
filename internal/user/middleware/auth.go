package middleware

import (
	"github.com/cemayan/event-scraper/config/user"
	"github.com/cemayan/event-scraper/internal/user/model"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// Protected protect routes
func Protected(configs *user.AppConfig) fiber.Handler {
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
