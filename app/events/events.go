package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

var EventsMap = make(map[string]Event)

const validPattern = "^[a-zA-Z0-9\u0400-\u04FF ]{3,50}$"

const (
	ErrDateFormat   = "неверный формат даты в событии: %s"
	ErrTitlePattern = "неверное имя в событии: %s"
	ErrorPriority   = "неверный приоритет в событии: %s"
	ErrorAddEvent   = "ошибка: %v в событии: %s"
)

type Priority string

const (
	PriorityHigh Priority = "High"
	PriorityLow  Priority = "Low"
)

type Event struct {
	ID       string
	Title    string
	StartAt  time.Time
	Priority Priority
}

func getNextID() string {
	return uuid.New().String()
}

func IsValidPriority(priority Priority) bool {
	if priority != PriorityHigh && priority != PriorityLow {
		return false
	}
	return true
}

func IsValidTitle(title string) bool {
	matched, err := regexp.MatchString(validPattern, title)
	if err != nil {
		return false
	}
	return matched
}

func IsValidDate(dateStr string) (time.Time, error) {
	t, err := dateparse.ParseAny(dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t,
		nil
}

func NewEvent(title string, dateStr string, priority Priority) (Event, error) {
	if !IsValidPriority(priority) {
		return Event{}, fmt.Errorf(ErrorPriority, title)
	}
	if !IsValidTitle(title) {
		return Event{}, fmt.Errorf(ErrTitlePattern, title)
	}
	t, err := IsValidDate(dateStr)
	if err != nil {
		return Event{}, fmt.Errorf(ErrorAddEvent, err, title)
	}
	return Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t,
		Priority: priority,
	}, nil
}
