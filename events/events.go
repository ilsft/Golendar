package events

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ilsft/Golendar/reminder"
	validators "github.com/ilsft/Golendar/utils"
)

const (
	errTitlePattern    = "неверное имя в событии: %s"
	errorValidEvent    = "ошибка %w в событии: %s"
	ErrorValidReminder = "ошибка %w в напоминании: %s"
)

type Event struct {
	ID       string             `json:"id"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder"`
}

func getNextID() string {
	return uuid.New().String()
}

func NewEvent(title string, dateStr string, priority Priority) (*Event, error) {
	err := validators.CheckTitleEmpty(title)
	if err != nil {
		return nil, err
	}
	if !validators.IsValidTitle(title) {
		return nil, fmt.Errorf(errTitlePattern, title)
	}
	t, err := validators.ValidateDate(dateStr)
	if err != nil {
		return nil, fmt.Errorf(errorValidEvent, err, title)
	}

	errPriority := priority.ValidatePriority()
	if errPriority != nil {
		return nil, fmt.Errorf(errorValidEvent, errPriority, title)
	}

	return &Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t,
		Priority: priority,
		Reminder: nil,
	}, nil
}

func (e *Event) Update(title string, date string, priority Priority) error {
	validEvent, err := NewEvent(title, date, priority)
	if err != nil {
		return err
	}
	e.Title = validEvent.Title
	e.StartAt = validEvent.StartAt
	return nil
}

func (e *Event) AddReminder(message string, at time.Time, notifier reminder.Notifier) (string, error) {
	rem, err := reminder.NewReminder(message, at, notifier)
	if err != nil {
		return "", fmt.Errorf(ErrorValidReminder, err, message)
	}

	e.Reminder = rem
	msg := e.Reminder.Start()
	return msg, nil
}

func (e *Event) RemoveReminder() string {
	msg := e.Reminder.Stop()
	e.Reminder = nil
	return msg
}
