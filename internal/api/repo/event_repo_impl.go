package repo

import (
	"github.com/cemayan/event-scraper/internal/api/utils"
	"github.com/cemayan/event-scraper/protos"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EventRepo struct {
	db  *gorm.DB
	log *log.Logger
}

// GetByProvider returns Events based on given provider
func (b EventRepo) GetByProvider(provider utils.Provider, page int, pageSize int) []protos.Event {
	var events []protos.Event
	p1, p2 := paginate(page, pageSize)
	if err := b.db.Where("provider = ? offset ? limit ? ", provider.String(), p1, p2).Find(&events); err.Error != nil {
		log.Warningln("provider not found!")
	}
	return events
}

// Paginate returns computed values based on given page and page_size
func paginate(page int, pageSize int) (int, int) {

	if page == 0 {
		page = 1
	}

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize == 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	return offset, pageSize
}

// DeleteByProvider removes based on given provider
func (b EventRepo) DeleteByProvider(provider string) {
	tx := b.db.Begin()
	tx.Where("provider = ?", provider).Delete(&protos.Event{})
	tx.Commit()
	b.log.Infof("%s provider events are deleted.", provider)
}

func (b EventRepo) Create(event *protos.Event) (*protos.Event, error) {
	tx := b.db.Begin()
	tx.Create(&event)
	tx.Commit()
	return event, nil
}

func NewEventRepo(db *gorm.DB, log *log.Logger) EventRepository {
	return &EventRepo{
		db:  db,
		log: log,
	}
}
