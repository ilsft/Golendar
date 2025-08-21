package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ilsft/app/events"
	"github.com/ilsft/app/storage"
	validators "github.com/ilsft/app/utils"
)

const (
	EventAddedMessage     = "Событие: %s добавлено"
	EventDeleteMessage    = "Событие: %s удалено"
	EventEditTitleMessage = "Событие: %s обновлено на %s - %s"
	ReminderAddMessage    = "Напоминание: %s добавлено \n%s"
	ReminderDeleteMessage = "Напоминание удалено \n%s"
	ReminderStopMessage   = "Напоминание: %s остановлено"
)

const (
	ErrorNotFoundID   = "не найдено событие с ID %s"
	ErrorEmptyList    = "событий нет"
	ErrorSerialJSON   = "ошибка сериализации: %v"
	ErrorDeSerialJSON = "ошибка десериализации: %v"
	ErrorNotFoundRem  = "нет напоминания для удаления"
)

type Calendar struct {
	CalendarEvents map[string]*events.Event `json:"events"`
	Storage        storage.Store            `json:"-"`
	Notification   chan string              `json:"-"`
}

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		CalendarEvents: make(map[string]*events.Event),
		Storage:        s,
		Notification:   make(chan string, 5),
	}
}

func (c *Calendar) AddEvent(title string, dateStr string, priority events.Priority) (string, error) {
	event, err := events.NewEvent(title, dateStr, priority)
	if err != nil {
		return "", err
	}
	c.CalendarEvents[event.ID] = event
	return fmt.Sprintf(EventAddedMessage, event.Title), nil
}

func (c *Calendar) ShowEvents() string {
	if len(c.CalendarEvents) == 0 {
		return ErrorEmptyList
	}
	var msgs []string
	for _, event := range c.CalendarEvents {
		msg := fmt.Sprintf("%s - %s - %v - %s", event.ID, event.Title,
			validators.FormatDateEvent(event.StartAt), event.Priority)
		msgs = append(msgs, msg)

		if event.Reminder != nil {
			remMsg := fmt.Sprintf("Напоминание для: %s - %s - %s - %v", event.Title,
				event.Reminder.Message, validators.FormatDateEvent(event.Reminder.At), event.Reminder.Sent)
			msgs = append(msgs, remMsg)
		}
	}
	return strings.Join(msgs, "\n")
}

func (c *Calendar) GetEventByID(id string) (*events.Event, error) {
	e, exist := c.CalendarEvents[id]
	if !exist {
		return nil, fmt.Errorf(ErrorNotFoundID, id)
	}
	return e, nil
}

func (c *Calendar) DeleteEvent(id string) (string, error) {
	event, err := c.GetEventByID(id)
	if err != nil {
		return "", err
	}
	delete(c.CalendarEvents, id)
	return fmt.Sprintf(EventDeleteMessage, event.Title), nil
}

func (c *Calendar) EditEvent(id string, newTitle string, date string, priority events.Priority) (string, error) {
	event, err := c.GetEventByID(id)
	if err != nil {
		return "", err
	}
	oldTitle := c.CalendarEvents[id].Title
	err = event.Update(newTitle, date, priority)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(EventEditTitleMessage, oldTitle, newTitle, date), nil
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf(ErrorSerialJSON, err)
	}
	err = c.Storage.Save(data)
	if err != nil {
		return (err)
	}
	return nil
}

func (c *Calendar) Load() error {
	data, err := c.Storage.Load()
	if err != nil {
		c.Notify(err.Error())
		return err
	}
	if len(data) == 0 {
		if c.CalendarEvents == nil {
			c.CalendarEvents = make(map[string]*events.Event)
		}
		return nil
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf(ErrorDeSerialJSON, err)
	}
	return nil
}

func (c *Calendar) SetEventReminder(id string, message string, time string) (string, error) {
	event, err := c.GetEventByID(id)
	if err != nil {
		return "", err
	}
	t, errValid := validators.ValidateDate(time)
	if errValid != nil {
		return "", fmt.Errorf(events.ErrorValidReminder, errValid, message)
	}
	msg, err := event.AddReminder(message, t, c)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(ReminderAddMessage, event.Reminder.Message, msg), nil
}

func (c *Calendar) RemoveEventReminder(id string) (string, error) {
	event, err := c.GetEventByID(id)
	if err != nil {
		return "", err
	}
	if event.Reminder == nil {
		return "", errors.New(ErrorNotFoundRem)
	}
	msg := event.RemoveReminder()
	return fmt.Sprintf(ReminderDeleteMessage, msg), nil
}

func (c *Calendar) CancelEventReminder(id string) (string, error) {
	event, err := c.GetEventByID(id)
	if err != nil {
		return "", err
	}
	msg := event.Reminder.Stop()
	return msg, nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}
