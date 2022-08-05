package repo

import (
	"github.com/cemayan/event-scraper-common/protos"
	"github.com/cemayan/event-scraper/api/src/utils"
)

type EventRepository interface {
	GetByProvider(provider utils.Provider, page int, pageSize int) []protos.Event
	Create(event *protos.Event) (*protos.Event, error)
	DeleteByProvider(provider string)
}
