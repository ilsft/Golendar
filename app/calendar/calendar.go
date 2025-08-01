package calendar

import (
	"fmt"

	"github.com/ilsft/app/events"
)

const timePattern = "Mon 2006/01/02 - 15:04"
const (
	EventAddedMessage  = "Событие: %s добавлено\n"
	EventDeleteMessage = "Событие: %s удалено\n"
	EventEditMessage   = "Событие: %s изменено\n"
	EventShowMessage   = "/////Cписок событий/////"
)

const (
	ErrorNotFoundDeleteEvent = "При удалении не найдено событие: %s \n"
	ErrorNotFoundEditEvent   = "При изменении не найдено событие: %s \n"
	ErrorParseMessage        = "%v в событии: %s\n"
	ErrAddEvent              = "%v при добавлении: %s\n"
)

var EventsMap = make(map[string]events.Event)

func AddEvent(title string, dateStr string, priority events.Priority) {
	event, err := events.NewEvent(title, dateStr, priority)
	if err != nil {
		fmt.Printf(ErrAddEvent, err, title)
		return
	}
	EventsMap[event.ID] = event
	fmt.Printf(EventAddedMessage, event.Title)
}

func ShowEvents() {
	fmt.Println(EventShowMessage)
	for _, event := range EventsMap {
		fmt.Printf("%s - %v \n", event.Title, event.StartAt.Format(timePattern))
	}
}

func searchID(title string) (string, bool) {
	for _, event := range EventsMap {
		if title == event.Title {
			return event.ID, true
		}
	}
	return "", false
}

func DeleteEvent(title string) {
	titleID, found := searchID(title)
	if !found {
		fmt.Printf(ErrorNotFoundDeleteEvent, title)
		return
	}
	delete(EventsMap, titleID)
	fmt.Printf(EventDeleteMessage, title)
}

func EditEvent(title string, dateStr string, priority events.Priority) {
	titleID, found := searchID(title)
	if !found {
		fmt.Printf(ErrorNotFoundEditEvent, title)
		return
	}
	t, err := events.IsValidDate(dateStr)
	if err != nil {
		fmt.Printf(ErrorParseMessage, err, title)
		return
	}
	if !events.IsValidPriority(priority) {
		fmt.Printf(ErrorParseMessage, events.ErrorPriority, title)
		return
	}

	e := events.Event{
		ID:       titleID,
		Title:    title,
		StartAt:  t,
		Priority: priority,
	}

	EventsMap[titleID] = e
	fmt.Printf(EventEditMessage, title)
}
