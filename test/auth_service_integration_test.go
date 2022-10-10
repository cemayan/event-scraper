package test

import (
	"bytes"
	"encoding/json"
	"github.com/cemayan/event-scraper/config/user"
	"github.com/cemayan/event-scraper/internal/user/database"
	database2 "github.com/cemayan/event-scraper/internal/user/database"
	"github.com/cemayan/event-scraper/internal/user/model"
	"github.com/cemayan/event-scraper/internal/user/repo"
	"github.com/cemayan/event-scraper/internal/user/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http/httptest"
	"os"
	"testing"
)

type e2eAuthTestSuite struct {
	suite.Suite
	app       *fiber.App
	db        *gorm.DB
	usrSvc    service.UserService
	authSvc   service.AuthService
	configs   *user.AppConfig
	v         *viper.Viper
	dbHandler database.DBHandler
}

func TestE2EAuthTestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (ts *e2eAuthTestSuite) SetupSuite() {
	app := fiber.New()
	app.Use(cors.New())

	ts.v = viper.New()
	_configs := user.NewConfig(ts.v)

	env := os.Getenv("ENV")
	appConfig, err := _configs.GetConfig(env)
	ts.configs = appConfig
	if err != nil {
		return
	}

	//Postresql connection
	ts.dbHandler = database.NewDbHandler(ts.configs)
	ts.dbHandler.ConnectDB()

	ts.app = app
	DB := database.DB
	ts.db = DB

	userRepo := repo.NewUserRepo(database2.DB, log.New())
	authSvc := service.NewAuthService(userRepo, log.New(), ts.configs)
	ts.authSvc = authSvc

	userSvc := service.NewUserService(userRepo, authSvc, log.New(), ts.configs)
	ts.usrSvc = userSvc

}

func (ts *e2eAuthTestSuite) removeAllRecords() {
	ts.db.Exec("DELETE FROM users")
}

func (ts *e2eAuthTestSuite) getRecords() []model.User {
	var users []model.User
	ts.db.Find(&users)
	return users
}

func (ts *e2eAuthTestSuite) getUserModel() model.User {
	var user model.User
	user.Password = "123"
	user.Username = "test"
	user.Email = "user@test.com"
	return user
}

func (ts *e2eAuthTestSuite) TestUserService_Create() {

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
