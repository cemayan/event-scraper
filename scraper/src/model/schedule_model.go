package model

import (
	"github.com/cemayan/event-scraper/scraper/src/scheduler"
	"github.com/cemayan/event-scraper/scraper/src/utils"
)

// ScheduleModel is a representation of based on given JSON config model
// CronJob is started with this params
type ScheduleModel struct {
	Name       string             `json:"name"`
	TimeType   scheduler.TimeType `json:"timeType"`
	TimeCount  uint64             `json:"timeCount"`
	At         *string            `json:"at"`
	Provider   utils.Provider     `json:"provider"`
	City       utils.City         `json:"city"`
	Category   utils.Category     `json:"category"`
	DatePeriod utils.DatePeriod   `json:"datePeriod"`
}
