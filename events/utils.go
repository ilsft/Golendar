package events

import (
	"errors"
	"time"
)

const ErrorNotFoundEvent = "событие не найдено"
const timePattern = "Mon 2006/01/02 - 15:04"

func FormatDateEvent(date time.Time) string {
	return date.Format(timePattern)
}

func SearchID(title string) (string, bool) {
	for _, event := range EventsMap {
		if title == event.Title {
			return event.ID, true
		}
	}
	return "", false
}

func FindEventIDByTitle(title string) (Event, error) {
	titleID, found := SearchID(title)
	if !found {
		return Event{}, errors.New(ErrorNotFoundEvent)
	}
	event := EventsMap[titleID]
	return event, nil
}
