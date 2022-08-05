package service

import (
	"bytes"
	"encoding/json"
	"github.com/cemayan/event-scraper/user/src/database"
	"github.com/cemayan/event-scraper/user/src/model"
	"github.com/cemayan/event-scraper/user/src/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	app     *fiber.App
	db      *gorm.DB
	usrSvc  UserService
	authSvc AuthService
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (ts *e2eTestSuite) SetupSuite() {
	app := fiber.New()
	app.Use(cors.New())

	ts.app = app
	DB := database.GetDB()
	ts.db = DB

	userRepo := repo.NewUserRepo(database.DB, log.New())
	authSvc := NewAuthService(userRepo, log.New())
	ts.authSvc = authSvc

	userSvc := NewUserService(userRepo, authSvc, log.New())
	ts.usrSvc = userSvc

}

func (ts *e2eTestSuite) removeAllRecords() {
	ts.db.Exec("DELETE FROM users")
}

func (ts *e2eTestSuite) getRecords() []model.User {
	var users []model.User
	ts.db.Find(&users)
	return users
}

func (ts *e2eTestSuite) getUserModel() model.User {
	var user model.User
	user.Password = "123"
	user.Username = "test"
	user.Email = "user@test.com"
	return user
}

func (ts *e2eTestSuite) TestUserService_Create() {

	ts.removeAllRecords()

	ts.app.Post("/user", ts.usrSvc.CreateUser)

	marshal, err := json.Marshal(ts.getUserModel())
	if err != nil {
		return
	}

	requestBuffer := bytes.NewBuffer(marshal)

	req := httptest.NewRequest("POST", "/user", requestBuffer)
	req.Header.Add("Content-Type", "application/json")

	_, err = ts.app.Test(req, 10000)
	if err != nil {
		return
	}

	users := ts.getRecords()
	ts.Equal(uint(1), users[0].ID)
	ts.Equal("test", users[0].Username)
	ts.Equal("user@test.com", users[0].Email)
}
