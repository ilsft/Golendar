package calendar

import (
	"fmt"

	"github.com/ilsft/app/events"
)

const (
	EventAddedMessage        = "Событие: %s добавлено"
	EventDeleteMessage       = "Событие: %s удалено"
	EventEditTitleMessage    = "Имя события: %s изменено на: %s"
	EventEditDateMessage     = "Дата события: %s изменена на: %v"
	EventEditPriorityMessage = "Приоритет события: %s изменена на: %v"
	EventShowMessage         = "/////Cписок событий/////"
)

const (
	ErrorNotFoundDeleteEvent = "при удалении не найдено событие: %s"
	ErrorParseMessage        = "%v в событии: %s"
	ErrAddEvent              = "%v при добавлении: %s"
)

func AddEvent(title string, dateStr string, priority events.Priority) (string, error) {
	event, err := events.NewEvent(title, dateStr, priority)
	if err != nil {
		return (""), err
	}
	events.EventsMap[event.ID] = event
	message := fmt.Sprintf(EventAddedMessage, event.Title)
	return message, nil
}

func ShowEvents() {
	fmt.Println(EventShowMessage)
	for _, event := range events.EventsMap {
		fmt.Printf("%s - %v - %s \n", event.Title, events.FormatDateEvent(event.StartAt), event.Priority)
	}
}

func DeleteEvent(title string) (string, error) {
	titleID, found := events.SearchID(title)
	if !found {
		return "", fmt.Errorf(ErrorNotFoundDeleteEvent, title)
	}
	delete(events.EventsMap, titleID)
	message := fmt.Sprintf(EventDeleteMessage, title)
	return message, nil
}

func EditTitleEvent(oldTitle string, newTitle string) (string, error) {
	event, err := events.FindEventIDByTitle(oldTitle)
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, oldTitle)
	}
	if !events.IsValidTitle(newTitle) {
		return "", fmt.Errorf(events.ErrTitlePattern, newTitle)
	}
	event.Title = newTitle
	events.EventsMap[event.ID] = event
	message := fmt.Sprintf(EventEditTitleMessage, oldTitle, newTitle)
	return message, nil
}

func EditDateEvent(title string, dateStr string) (string, error) {
	event, err := events.FindEventIDByTitle(title)
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, title)
	}
	t, err := events.IsValidDate(dateStr)
	if err != nil {
		return "", fmt.Errorf(ErrorParseMessage, err, title)
	}
	event.StartAt = t
	events.EventsMap[event.ID] = event
	message := fmt.Sprintf(EventEditDateMessage, title, events.FormatDateEvent(t))
	return message, nil
}

func EditPriorityEvent(title string, priority events.Priority) (string, error) {
	event, err := events.FindEventIDByTitle(title)
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, title)
	}
	if !events.IsValidPriority(priority) {
		return "", fmt.Errorf(events.ErrorPriority, title)
	}
	event.Priority = priority
	events.EventsMap[event.ID] = event
	message := fmt.Sprintf(EventEditPriorityMessage, title, priority)
	return message, nil
}
