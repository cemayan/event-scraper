package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Provider int

func (p Provider) String() string {
	return [...]string{"BILETIX", "PASSO", "KULTURIST"}[p]
}

const (
	BILETIX Provider = iota
	PASSO
	KULTURIST
)

// FailOnError returns a log based on given error and message
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// GetTokenFromHeaders return based on given Authorization token in http Request
func GetTokenFromHeaders(c *fiber.Ctx) (string, error) {

	authHeader := c.GetReqHeaders()["Authorization"]

	authArr := strings.Split(authHeader, "Bearer ")
	if len(authArr) != 2 {
		return "", fmt.Errorf("invalid Authorization header")
	}
	return authArr[1], nil
}
