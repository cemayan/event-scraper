package test

import (
	"encoding/json"
	"github.com/cemayan/event-scraper/config/api"
	"github.com/cemayan/event-scraper/internal/api/database"
	"github.com/cemayan/event-scraper/internal/api/repo"
	"github.com/cemayan/event-scraper/internal/api/service"
	"github.com/cemayan/event-scraper/protos"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
)

type e2eTestSuite struct {
	suite.Suite
	app       *fiber.App
	db        *gorm.DB
	eSvc      service.EventService
	configs   *api.AppConfig
	v         *viper.Viper
	dbHandler database.DBHandler
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (ts *e2eTestSuite) SetupSuite() {
	app := fiber.New()
	app.Use(cors.New())

	ts.v = viper.New()
	_configs := api.NewConfig(ts.v)

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

	eRepo := repo.NewEventRepo(DB, log.New())
	ts.eSvc = service.NewEventService(eRepo, log.New())
}

func (ts *e2eTestSuite) createSomeRecord() {
	eventModel := protos.Event{

		Type:       "MUSIC",
		EventName:  "TEST_EVENT",
		Place:      "TEST_PLACE",
		FirstDate:  "2022-08-17 18:00:00 +0000 UTC",
		SecondDate: "2022-08-17 18:00:00 +0000 UTC",
		Provider:   "BILETIX",
	}

	eventModel2 := protos.Event{
		Type:       "ART",
		EventName:  "TEST_EVENT2",
		Place:      "TEST_PLACE2",
		FirstDate:  "2022-08-17 18:00:00 +0000 UTC",
		SecondDate: "2022-08-17 18:00:00 +0000 UTC",
		Provider:   "BILETIX",
	}

	eventModel3 := protos.Event{
		Type:       "MUSIC",
		EventName:  "TEST_EVENT3",
		Place:      "TEST_PLACE3",
		FirstDate:  "2022-08-17 18:00:00 +0000 UTC",
		SecondDate: "2022-08-17 18:00:00 +0000 UTC",
		Provider:   "PASSO",
	}

	ts.db.Create(&eventModel)
	ts.db.Create(&eventModel2)
	ts.db.Create(&eventModel3)
}

func (ts *e2eTestSuite) removeAllRecords() {
	ts.db.Exec("DELETE FROM events")
}

func (ts *e2eTestSuite) getRecords() []protos.Event {
	var events []protos.Event
	ts.db.Find(&events)
	return events
}

func (ts *e2eTestSuite) TestEventService_HealthCheck() {
	ts.removeAllRecords()
	ts.createSomeRecord()

	ts.app.Get("/health", ts.eSvc.HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Add("Content-Type", "application/json")

	test, err := ts.app.Test(req, 10)
	if err != nil {
		return
	}
	ts.Equal(200, test.StatusCode)
}

func (ts *e2eTestSuite) TestEventService_GetByProvider() {
	ts.removeAllRecords()
	ts.createSomeRecord()

	ts.app.Get("/event/provider/:provider", ts.eSvc.GetByProvider)

	req := httptest.NewRequest("GET", "/event/provider/0", nil)
	req.Header.Add("Content-Type", "application/json")

	test, err := ts.app.Test(req, 10)
	if err != nil {
		return
	}
	ts.Equal(200, test.StatusCode)
	body, err := ioutil.ReadAll(test.Body)

	var events []protos.Event

	json.Unmarshal(body, &events)

	ts.Equal(2, len(events))
}

func (ts *e2eTestSuite) TestEventService_GetByProviderWithParams() {
	ts.removeAllRecords()
	ts.createSomeRecord()

	ts.app.Get("/event/provider/:provider", ts.eSvc.GetByProvider)

	req := httptest.NewRequest("GET", "/event/provider/0?page=1&page_size=1", nil)
	req.Header.Add("Content-Type", "application/json")

	test, err := ts.app.Test(req, 10)
	if err != nil {
		return
	}
	ts.Equal(200, test.StatusCode)
	body, err := ioutil.ReadAll(test.Body)

	var events []protos.Event

	json.Unmarshal(body, &events)

	ts.Equal(1, len(events))
}

func (ts *e2eTestSuite) TestEventService_DeleteByProvider() {
	ts.removeAllRecords()
	ts.createSomeRecord()

	events := ts.getRecords()
	ts.Equal(3, len(events))

	ts.eSvc.DeleteByProvider("PASSO")

	events2 := ts.getRecords()

	ts.Equal(2, len(events2))
}
