package events

import (
	"errors"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
)

const validPattern = "^[a-zA-Z0-9\u0400-\u04FF ]{3,50}$"

const (
	ErrDateFormat    = "неверный формат даты"
	ErrTitlePatternt = "неверное имя события"
	ErrorPriority    = "неверный приоритет"
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
	pattern := validPattern
	matched, err := regexp.MatchString(pattern, title)
	if err != nil {
		return false
	}
	return matched
}

func IsValidDate(dateStr string) (time.Time, error) {
	t, err := dateparse.ParseAny(dateStr)
	if err != nil {
		return time.Time{}, errors.New(ErrDateFormat)
	}
	return t,
		nil
}

func NewEvent(title string, dateStr string, priority Priority) (Event, error) {
	if !IsValidPriority(priority) {
		return Event{}, errors.New(ErrorPriority)
	}
	if !IsValidTitle(title) {
		return Event{}, errors.New(ErrTitlePatternt)
	}
	t, err := IsValidDate(dateStr)
	if err != nil {
		return Event{}, err
	}
	return Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t,
		Priority: priority,
	}, nil
}
