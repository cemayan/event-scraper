package service

import (
	"fmt"
	"github.com/cemayan/event-scraper/config/user"
	"github.com/cemayan/event-scraper/internal/user/model"
	"github.com/cemayan/event-scraper/internal/user/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
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

// A AuthSvc  contains the required dependencies for this service
type AuthSvc struct {
	repository model.UserRepository
	log        *log.Logger
	configs    *user.AppConfig
}

// checkPasswordHash  returns valid status based on given password
func (s AuthSvc) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s AuthSvc) isValidUser(id uint) bool {

	user, err := s.repository.GetUserById(id)
	if user == nil || err != nil {
		return false
	}
	return true
}

// validToken  returns valid status based on given token
// If based on given token is not found in claims map, it is returned false
func (s AuthSvc) isValidUserId(t *jwt.Token, id int) bool {

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	if uid != id {
		return false
	}

	return true
}

func (s AuthSvc) isValidUserWithPass(id uint, p string) bool {

	user, err := s.repository.GetUserById(id)
	if user == nil || err != nil {
		return false
	}

	if !s.checkPasswordHash(p, user.Password) {
		return false
	}
	return true
}

func (s AuthSvc) ValidToken(c *fiber.Ctx) error {

	token, err := utils.GetTokenFromHeaders(c, s.configs)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: fmt.Sprintf("JWT token parse error %v", err),
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	uid := uint(claims["user_id"].(float64))
	if isValid := s.isValidUser(uid); !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: fmt.Sprintf("Invalid user %v", uid),
		})
	}

	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{Message: "invalid token!"})
	}
	return c.Status(200).JSON(model.Response{Message: "valid token!"})
}

// HealthCheck returns 200 with body
func (s AuthSvc) HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON("healty!")
}

// CheckPasswordHash returns the correctness of password
func (s AuthSvc) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CheckPasswordHash returns user  based on given email.
func (s AuthSvc) getUserByEmail(e string) (*model.User, error) {
	return s.repository.GetUserByEmail(e)
}

// CheckPasswordHash returns user  based on given username.
func (s AuthSvc) getUserByUsername(u string) (*model.User, error) {
	return s.repository.GetUserByUsername(u)
}

// Login returns authentication result
// If given password or username is not correct, it is returned 403
// Then, it is created new jwt token. Username,user_id and exp is added to token claims.
func (s AuthSvc) Login(c *fiber.Ctx) error {

	var input model.LoginInput
	var ud model.UserData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: fmt.Sprintf("Error on login request %v", err),
		})
	}
	identity := input.Username
	pass := input.Password

	user, err := s.getUserByUsername(identity)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Message: fmt.Sprintf("Error on username %v", err),
		})
	}

	ud = model.UserData{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	if !s.CheckPasswordHash(pass, ud.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Message: fmt.Sprintf("Invalid password %v", nil),
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["user_id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	t, err := token.SignedString([]byte(s.configs.SECRET))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Message: fmt.Sprintf("Invalid  secret %v", nil),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{
		Message: fmt.Sprintf("%v", t),
	})

}

func NewAuthService(rep model.UserRepository, log *log.Logger, configs *user.AppConfig) AuthService {
	return &AuthSvc{
		repository: rep,
		log:        log,
		configs:    configs,
	}
}
