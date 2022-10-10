package common

import (
	"github.com/google/uuid"
)

type EventName string

const (
	DELETE_EVENTS_IN_TABLE EventName = "DELETE_EVENTS_IN_TABLE"
)

type ScraperEvent struct {
	AggregationId uuid.UUID
	EventName     EventName
	EventDate     int64
	Payload       interface{}
}
