package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/shlex"
	"github.com/ilsft/Golendar/events"
)

const (
	errIncorrectChoice = "неверный выбор"
	errLenEmptyTitle   = "название события не указано"
	errNoMatchTitle    = "совпадений не найдено"
)

func (c *Cmd) readLineWithPrompt(prompt string) (string, error) {
	c.handlePrint(prompt)
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func (c *Cmd) getUserChoice(max int) (int, error) {
	line, err := c.readLineWithPrompt(inputNumberMessage)
	if err != nil {
		return 0, err
	}
	choice, err := strconv.Atoi(line)
	if err != nil || choice < 1 || choice > max {
		return 0, errors.New(errIncorrectChoice)
	}
	return choice, nil
}

func (c *Cmd) readAndParseInput(promptMsg string) ([]string, error) {
	line, err := c.readLineWithPrompt(promptMsg)
	if err != nil {
		return nil, err
	}
	return shlex.Split(line)
}

func (c *Cmd) titleMatches(title string, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(title), strings.ToLower(prefix))
}

func (c *Cmd) selectEvents(parts []string) (*events.Event, error) {
	if len(parts) <= 1 {
		return nil, errors.New(errLenEmptyTitle)
	}
	var matchedEvents []*events.Event
	for _, event := range c.calendar.CalendarEvents {
		if c.titleMatches(event.Title, parts[1]) {
			matchedEvents = append(matchedEvents, event)
		}
	}
	return c.chooseEvent(matchedEvents)
}

func (c *Cmd) selectEventsByReminder(showWithReminders bool, parts []string) (*events.Event, error) {
	if len(parts) <= 1 {
		return nil, errors.New(errLenEmptyTitle)
	}
	var matchedEvents []*events.Event
	prefix := parts[1]
	for _, event := range c.calendar.CalendarEvents {
		hasReminder := event.Reminder != nil && event.Reminder.Message != ""
		if c.titleMatches(event.Title, prefix) {
			if showWithReminders && hasReminder {
				matchedEvents = append(matchedEvents, event)
			} else if !showWithReminders && !hasReminder {
				matchedEvents = append(matchedEvents, event)
			}
		}
	}
	return c.chooseEvent(matchedEvents)
}

func (c *Cmd) chooseEvent(matchedEvents []*events.Event) (*events.Event, error) {
	if len(matchedEvents) == 0 {
		return nil, errors.New(errNoMatchTitle)
	}
	for i, event := range matchedEvents {
		fmt.Printf("%d. %s - %s\n", i+1, event.Title, event.StartAt.Format(patternTime))
	}
	choice, err := c.getUserChoice(len(matchedEvents))
	if err != nil {
		return nil, err
	}
	return matchedEvents[choice-1], nil
}
