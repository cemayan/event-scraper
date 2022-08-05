package service

import (
	"github.com/cemayan/event-scraper/user/src/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	checkPasswordHash(password, hash string) bool
	isValidUser(id uint) bool
	isValidUserId(t *jwt.Token, id int) bool
	isValidUserWithPass(id uint, p string) bool
	HealthCheck(c *fiber.Ctx) error
	CheckPasswordHash(password, hash string) bool
	getUserByEmail(e string) (*model.User, error)
	getUserByUsername(u string) (*model.User, error)
	Login(c *fiber.Ctx) error
	ValidToken(c *fiber.Ctx) error
}
