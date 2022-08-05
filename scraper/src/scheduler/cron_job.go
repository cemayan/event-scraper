package scheduler

import "github.com/jasonlvhit/gocron"

type TimeType int

const (
	SECOND TimeType = iota
	MINUTE
	HOUR
	DAY
	WEEK
)

type SchedulerService interface {
	Start() bool
}

// SchedulerSvc is representation of a dependencies
type SchedulerSvc struct {
	_scheduler *gocron.Scheduler
}

// SetJob is used to set Job based on given payload
// "timeCount" refers to interval of time
// "at" refers to specific time such as "10:30:01"
// "T" refers to which handler will execute
// "params refers to given handler parameters
func SetJob[T any](_scheduler *gocron.Scheduler,
	timeType TimeType,
	timeCount uint64,
	at *string,
	t T, params ...interface{}) {
	switch timeType {
	case SECOND:
		if timeCount == uint64(1) {
			_scheduler.Every(1).Second().Do(t, params)
		} else {
			_scheduler.Every(timeCount).Seconds().Do(t, params)
		}
	case MINUTE:
		if timeCount == uint64(1) {
			_scheduler.Every(1).Minute().Do(t, params)
		} else {
			_scheduler.Every(timeCount).Minutes().Do(t, params)
		}
	case HOUR:
		if timeCount == uint64(1) {
			_scheduler.Every(1).Hour().Do(t, params)
		} else {
			_scheduler.Every(timeCount).Hour().Do(t, params)
		}
	case DAY:
		if at != nil {
			_scheduler.Every(1).Day().At(*at).Do(t, params)
		} else {
			_scheduler.Every(1).Day().Do(t, params)
		}
	}
}

// Start starts the scheduler
func (s SchedulerSvc) Start() bool {
	return <-s._scheduler.Start()
}

func NewScheduler(_scheduler *gocron.Scheduler) SchedulerService {
	return &SchedulerSvc{_scheduler: _scheduler}
}
