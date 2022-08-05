package middleware

import (
	"github.com/cemayan/event-scraper/api/src/config"
	"github.com/cemayan/event-scraper/api/src/model"
	"github.com/cemayan/event-scraper/api/src/utils"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	log "github.com/sirupsen/logrus"
	"github.com/sony/gobreaker"
)

type Middleware interface {
	Protected() fiber.Handler
	ValidUserId(c *fiber.Ctx) error
}

type MiddlewareSvc struct {
	httpClient *resty.Client
	log        *log.Logger
	cb         *gobreaker.CircuitBreaker
}

// ValidUserId checks header token  based on given token
// Token user id  might be different than desired user
func (m MiddlewareSvc) ValidUserId(c *fiber.Ctx) error {

	configs := config.GetConfig()
	token, _ := utils.GetTokenFromHeaders(c)

	statusCode, err := m.cb.Execute(func() (interface{}, error) {
		resp, err := m.httpClient.R().
			EnableTrace().
			SetAuthToken(token).
			Post("http://" + configs.AUTH_SERVER + "/api/v1/auth/validateToken")

		if err != nil {
			return nil, err
		}

		return resp.StatusCode(), nil
	})

	if statusCode != fiber.StatusOK {
		return c.Status(fiber.StatusUnauthorized).
			JSON(model.Response{Message: "Invalid request"})
	}

	err = c.Next()
	if err != nil {
		return err
	}
	return err
}

// Protected protect routes
func (m MiddlewareSvc) Protected() fiber.Handler {

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

func NewMiddleware(httpClient *resty.Client, cb *gobreaker.CircuitBreaker, log *log.Logger) Middleware {
	return &MiddlewareSvc{httpClient: httpClient, cb: cb, log: log}
}
