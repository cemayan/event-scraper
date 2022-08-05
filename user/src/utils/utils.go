package utils

import (
	"fmt"
	"github.com/cemayan/event-scraper/user/src/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"strings"
)

// FailOnError returns a log based on given error and message
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func GetTokenFromHeaders(c *fiber.Ctx) (*jwt.Token, error) {
	configs := config.GetConfig()

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
