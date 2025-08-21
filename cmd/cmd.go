package cmd

import (
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/ilsft/Golendar/calendar"
	"github.com/ilsft/Golendar/logger"
)

type Cmd struct {
	calendar *calendar.Calendar
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
	}
}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "update", Description: "Изменить событие"},
		{Text: "add_rm", Description: "Добавить напоминание"},
		{Text: "stop_rm", Description: "Остановить напоминание"},
		{Text: "remove_rm", Description: "Удалить напоминание"},
		{Text: "log", Description: "Показать логи"},
		{Text: "help", Description: "Показать справку"},
		{Text: "exit", Description: "Выйти из программы"},
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) executor(input string) {
	parts, err := shlex.Split(input)
	if err != nil {
		c.handlePrint(err.Error())
		logger.LogError(err.Error())
		return
	}
	l.logMessage(input)
	if len(parts) == 0 {
		c.handlePrint(emptyInput)
		logger.LogError(emptyInput)
		return
	}
	cmd := strings.ToLower(parts[0])
	switch cmd {
	case "add":
		c.handleAddCmd(parts)
	case "remove":
		c.handleDeleteCmd()
	case "update":
		c.handleEditeCmd(parts)
	case "add_rm":
		c.handleAddReminderCmd(parts)
	case "stop_rm":
		c.handleStopReminderCmd()
	case "remove_rm":
		c.handleDeleteReminderCmd()
	case "list":
		c.handleShowEventsCmd()
	case "log":
		c.handleShowLogsCmd()
	case "help":
		c.handleShowHelpCmd()
	case "exit":
		c.handleExitCmd()
	default:
		c.handleDefaultCmd(cmd)
	}

	err = c.calendar.Save()
	if err != nil {
		c.handlePrint(err.Error())
		return
	}
	err = l.saveLogs()
	if err != nil {
		c.handlePrint(err.Error())
		return
	}

}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	err := l.loadLogs()
	if err != nil {
		c.handlePrint(err.Error())

	}
	go func() {
		for msg := range c.calendar.Notification {
			c.handlePrint(msg)
		}
	}()
	p.Run()
}
