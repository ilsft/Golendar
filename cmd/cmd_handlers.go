package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/c-bata/go-prompt"
	"github.com/ilsft/Golendar/events"
	"github.com/ilsft/Golendar/logger"
	validators "github.com/ilsft/Golendar/utils"
)

const (
	errLenAddMessage         = "Формат: add \"название события\" \"дата и время\" \"приоритет\""
	errLenUpdMessage         = "Формат: update \"id события\" \"название события\" \"дата и время\" \"приоритет\""
	errLenAddReminderMessage = "Формат: add_rm \"id события\" \"название напоминания\" \"дата и время\""
)

const (
	unknownCommand        = "Неизвестная команда:"
	enterEventTitlePrompt = "Введите название события:"
	deafaultMessage       = "Введите 'help' для списка команд"
	emptyInput            = "Пустой ввод, повторите попытку"
)

const (
	ErrEmptyTitle     = "Запрещено создавать пустые события"
	ErrPastTimeTravel = "Запрещено путешествовать в прошлое"
)

const EventShowMessage = "📅Cписок событий✅"
const helpMessage = `add ✅ — добавить событие, ` + errLenAddMessage + `
remove ❌ — удалить событие, ввести название события вручную или выбрать из списка через <tab>
list 📒 — вывести все события в формате id - название - дата и время - приоритет
update ✏️ — изменить событие, найти id события по команде <list>, 
далее ввести данные, ` + errLenUpdMessage + `
add_rm 🔔 — добавить напоминание
stop_rm ⏸️ — остановить напоминание
remove_rm 🗑️ — удалить напоминание
log 📜 — показать логи
exit 🏁 — выход`

var idRegexp = regexp.MustCompile(`ID:([\w-]+)`)

func (c *Cmd) completerEvent(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{}
	for _, event := range c.calendar.CalendarEvents {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        fmt.Sprintf("%s | ID:%s", event.Title, event.ID[:4]),
			Description: fmt.Sprintf("%s | %s ", event.StartAt.Format("2006-01-02"), event.Priority),
		})
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) searchID(input string) string {
	matches := idRegexp.FindStringSubmatch(input)
	if len(matches) > 1 {
		shortID := matches[1]
		for _, event := range c.calendar.CalendarEvents {
			if len(event.ID) > 4 && event.ID[:4] == shortID {
				return event.ID
			}
		}
	}
	return (input)
}

func (c *Cmd) handleEventPrompt(f func(string) (string, error)) (string, error) {
	t := prompt.Input(enterEventTitlePrompt, c.completerEvent)
	fullEventID := c.searchID(t)
	msg, err := f(fullEventID)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (c *Cmd) notifyResult(msg string, err error) bool {
	if err != nil {
		c.handlePrint(err.Error())
		logger.LogError(err.Error())
		return false
	}
	c.handlePrint(msg)
	logger.LogInfo(msg)
	return true
}

func (c *Cmd) handleAddCmd(parts []string) {
	if len(parts) < 4 {
		c.handlePrint(errLenAddMessage)
		logger.LogError(errLenAddMessage)
		return
	}
	title := parts[1]
	date := parts[2]
	priority := events.Priority(parts[3])
	msg, err := c.calendar.AddEvent(title, date, priority)
	switch {
	case errors.Is(err, validators.ErrEmptyTitle):
		c.handlePrint(ErrEmptyTitle)
		logger.LogError(ErrEmptyTitle)
	case errors.Is(err, validators.ErrDateAlreadyPassed):
		c.handlePrint(ErrPastTimeTravel)
		logger.LogError(ErrPastTimeTravel)
	default:
		if !c.notifyResult(msg, err) {
			return
		}
	}
}
func (c *Cmd) handleDeleteCmd() {
	msg, err := c.handleEventPrompt(c.calendar.DeleteEvent)
	if !c.notifyResult(msg, err) {
		return
	}
}
func (c *Cmd) handleEditeCmd(parts []string) {
	if len(parts) < 5 {
		c.handlePrint(errLenUpdMessage)
		logger.LogError(errLenUpdMessage)
		return
	}
	id := parts[1]
	title := parts[2]
	date := parts[3]
	priority := events.Priority(parts[4])
	msg, err := c.calendar.EditEvent(id, title, date, priority)
	if !c.notifyResult(msg, err) {
		return
	}
}
func (c *Cmd) handleAddReminderCmd(parts []string) {
	if len(parts) < 4 {
		c.handlePrint(errLenAddReminderMessage)
		logger.LogError(errLenAddReminderMessage)
		return
	}
	id := parts[1]
	message := parts[2]
	time := parts[3]
	msg, err := c.calendar.SetEventReminder(id, message, time)
	if !c.notifyResult(msg, err) {
		return
	}
}
func (c *Cmd) handleStopReminderCmd() {
	msg, err := c.handleEventPrompt(c.calendar.CancelEventReminder)
	if !c.notifyResult(msg, err) {
		return
	}
}
func (c *Cmd) handleDeleteReminderCmd() {
	msg, err := c.handleEventPrompt(c.calendar.RemoveEventReminder)
	if !c.notifyResult(msg, err) {
		return
	}
}

func (c *Cmd) handleShowEventsCmd() {
	c.handlePrint(EventShowMessage)
	c.calendar.Notify(c.calendar.ShowEvents())
}

func (c *Cmd) handleShowLogsCmd() {
	fmt.Println(l.showLogs())
}
func (c *Cmd) handleShowHelpCmd() {
	c.handlePrint(helpMessage)
}

func (c *Cmd) handlePrint(msg string) {
	fmt.Println(msg)
	l.logMessage(msg)
}

func (c *Cmd) handleExitCmd() {
	close(c.calendar.Notification)
	os.Exit(0)
}

func (c *Cmd) handleDefaultCmd(cmd string) {
	msg := (unknownCommand + cmd)
	c.handlePrint(msg)
	c.handlePrint(deafaultMessage)
	logger.LogError(msg)
}
