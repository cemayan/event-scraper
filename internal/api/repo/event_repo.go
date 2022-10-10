package repo

import (
	"github.com/cemayan/event-scraper/internal/api/utils"
	"github.com/cemayan/event-scraper/protos"
)

type EventRepository interface {
	GetByProvider(provider utils.Provider, page int, pageSize int) []protos.Event
	Create(event *protos.Event) (*protos.Event, error)
	DeleteByProvider(provider string)
}
