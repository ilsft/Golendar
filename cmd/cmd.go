package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/ilsft/Golendar/calendar"
	"github.com/ilsft/Golendar/logger"
)

type Cmd struct {
	calendar *calendar.Calendar
	logger   *HistoryLogger
	reader   *bufio.Reader
}

func NewCmd(c *calendar.Calendar, logger *HistoryLogger) *Cmd {
	return &Cmd{
		calendar: c,
		logger:   logger,
		reader:   bufio.NewReader(os.Stdin),
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
		{Text: "history", Description: "Показать историю ввода/вывода"},
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

	c.logger.logMessage(input)
	if len(parts) == 0 {
		c.handlePrint(fmt.Sprint(emptyInput, "\n", deafaultMessage))
		logger.LogError(fmt.Sprint(emptyInput, "\n", deafaultMessage))
		return
	}
	cmd := strings.ToLower(parts[0])
	switch cmd {
	case "add":
		c.handleAddCmd(parts)
	case "remove":
		c.handleDeleteCmd(parts)
	case "update":
		c.handleEditeCmd(parts)
	case "add_rm":
		c.handleAddReminderCmd(parts)
	case "stop_rm":
		c.handleStopReminderCmd(parts)
	case "remove_rm":
		c.handleDeleteReminderCmd(parts)
	case "list":
		c.handleShowEventsCmd()
	case "history":
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
		logger.LogError(err.Error())
		return
	}
	err = c.logger.saveLogs()
	if err != nil {
		c.handlePrint(err.Error())
		logger.LogError(err.Error())
		return
	}

}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	err := c.logger.loadLogs()
	if err != nil {
		c.handlePrint(err.Error())
		logger.LogError(err.Error())
	}
	go func() {
		for msg := range c.calendar.Notification {
			c.handlePrint(msg)
		}
	}()
	p.Run()
}
