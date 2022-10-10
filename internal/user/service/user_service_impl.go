package service

import (
	"fmt"
	"github.com/cemayan/event-scraper/config/user"
	dto2 "github.com/cemayan/event-scraper/internal/user/dto"
	model2 "github.com/cemayan/event-scraper/internal/user/model"
	"github.com/cemayan/event-scraper/internal/user/utils"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	hashPassword(password string) (string, error)
	GetUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	HealthCheck(c *fiber.Ctx) error
}

// A UserSvc  contains the required dependencies for this service
type UserSvc struct {
	repository model2.UserRepository
	authSvc    AuthService
	log        *log.Logger
	configs    *user.AppConfig
}

func NewUserService(rep model2.UserRepository, authSvc AuthService, log *log.Logger, configs *user.AppConfig) UserService {
	return &UserSvc{
		repository: rep,
		authSvc:    authSvc,
		log:        log,
		configs:    configs,
	}
}

// HealthCheck returns 200 with body
func (s UserSvc) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON("healty!")
}

// hashPassword returns encrypted password based on given password
func (s UserSvc) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GetUser returns user based on given id.
func (s UserSvc) GetUser(c *fiber.Ctx) error {

	id, _ := c.ParamsInt("id")
	user, err := s.repository.GetUserById(uint(id))
	if user == nil || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model2.Response{
			Message: fmt.Sprintf("No user found with %v", user.ID),
		})
	}

	return c.JSON(model2.Response{
		Data: model2.UserData{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

// CreateUser creates new user based on given payload
// While user is creating password is encrypted then it is assigned as a password
func (s UserSvc) CreateUser(c *fiber.Ctx) error {
	user := new(model2.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model2.Response{
			Message: fmt.Sprintf("Review your input %s", err),
		})
	}

	hash, err := s.hashPassword(user.Password)
	fmt.Println(err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model2.Response{
			Message: fmt.Sprintf("Couldn't hash password %s", err),
		})
	}

	user.Password = hash
	_user, err := s.repository.CreateUser(user)
	if _user == nil || err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model2.Response{
			Message: fmt.Sprintf("Couldn't create use %s", err),
		})
	}

	newUser := dto2.NewUser{
		Email:    _user.Email,
		Username: _user.Username,
	}

	return c.Status(fiber.StatusCreated).JSON(model2.Response{
		Message: fmt.Sprintf("User created %s", newUser),
	})
}

// UpdateUser return updated user based on given payload
func (s UserSvc) UpdateUser(c *fiber.Ctx) error {

	token, _ := utils.GetTokenFromHeaders(c, s.configs)

	var userDTO dto2.UpdateUser
	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model2.Response{
			Message: fmt.Sprintf("Review your input %s", err),
		})
	}

	id, _ := c.ParamsInt("id")
	if id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model2.Response{
			Message: fmt.Sprintf("Review your id %v", id),
		})
	}

	isValidUser := s.authSvc.isValidUserId(token, id)

	if !isValidUser {
		return c.Status(fiber.StatusUnauthorized).JSON(model2.Response{
			Message: fmt.Sprintf("Given token has invalid user id"),
		})
	}

	_user, err := s.repository.GetUserById(uint(id))
	if _user == nil || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model2.Response{
			Message: fmt.Sprintf("No user found with %s", err),
		})
	}

	if userDTO.Password != "" {
		hash, _ := s.hashPassword(userDTO.Password)
		_user.Password = hash
	}
	if userDTO.Username != "" {
		_user.Username = userDTO.Username
	}
	if userDTO.Email != "" {
		_user.Email = userDTO.Email
	}

	s.repository.UpdateUser(_user)
	return c.Status(fiber.StatusOK).JSON(model2.Response{
		Message: fmt.Sprintf("User successfully updated"),
	})
}

// DeleteUser removes  the user based on given payload
func (s UserSvc) DeleteUser(c *fiber.Ctx) error {

	token, _ := utils.GetTokenFromHeaders(c, s.configs)

	id, _ := c.ParamsInt("id")

	isValidUser := s.authSvc.isValidUserId(token, id)
	if !isValidUser {
		return c.Status(fiber.StatusUnauthorized).JSON(model2.Response{
			Message: fmt.Sprintf("Given token has invalid user id"),
		})
	}

	s.repository.DeleteUser(uint(id))

	return c.Status(fiber.StatusOK).JSON(model2.Response{
		Message: fmt.Sprintf("User successfully deleted %v", id),
	})

}
