package utils

import (
	"fmt"
	"github.com/cemayan/event-scraper/config/user"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"strings"
)

// FailOnError returns a log based on given error and message
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetTokenFromHeaders(c *fiber.Ctx, configs *user.AppConfig) (*jwt.Token, error) {

	authHeader := c.GetReqHeaders()["Authorization"]

	authArr := strings.Split(authHeader, "Bearer ")
	if len(authArr) != 2 {
		return nil, fmt.Errorf("invalid Authorization header")
	}

	token, err := jwt.Parse(authArr[1], func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.SECRET), nil
	})

	return token, err
}
