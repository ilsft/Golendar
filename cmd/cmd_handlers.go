package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilsft/Golendar/events"
	"github.com/ilsft/Golendar/logger"
	validators "github.com/ilsft/Golendar/utils"
)

const (
	unknownCommand     = "Неизвестная команда:"
	deafaultMessage    = "Введите 'help' для списка команд"
	emptyInput         = "Пустой ввод, повторите попытку"
	inputNumberMessage = "Введите номер события: "
)

var (
	errEmptyTitle     = "запрещено создавать пустые события"
	errPastTimeTravel = "запрещено путешествовать в прошлое"
)

const patternTime = "2006-01-02 15:04:05"

const eventShowMessage = "📅Cписок событий✅"

const (
	errAddFormat      = `add "имя события" "дата и время" "приоритет"`
	errUpdateFormat   = `введите: "новое имя события" "новая дата и время" "новый приоритет"`
	errReminderFormat = `введите: "имя напоминания" "дата и время"`
)

const helpMessage = `
╔═════════════════════════════════════════════════╗
║                Справка по командам              ║
╚═════════════════════════════════════════════════╝

───────────[ Создание и просмотр событий ]───────────
  add      ✅    ┆ создать событие
                 ┆ формат: ` + errAddFormat + `
  list     📒    ┆ список всех событий 
                 ┆ (id - имя события - дата и время - приоритет)

────────────[ Работа с существующими событиями ]──────
  remove    ❌  ┆ удалить событие
  update    ✏️   ┆ изменить данные
                ┆ формат: ` + errUpdateFormat + `
  add_rm    🔔  ┆ добавить напоминание
                ┆ формат: ` + errReminderFormat + `
  stop_rm   ⏸️   ┆ остановить напоминание
  remove_rm 🗑️   ┆ удалить напоминание
		
──────────────[ Сервисные команды ]───────────────
  history   📜   ┆ показать журнал действий
  exit      🏁   ┆ выход из программы


═══════════[ Как работают команды для событий ]═══════════
1. Введите часть имени события:
       пример → "meet"
2. Программа покажет список совпадений:

   ┌───┬───────────────────────────────┬───────────────────┐
   │ № │ Имя события                   │ Дата и время      │
   ├───┼───────────────────────────────┼───────────────────┤
   │ 1 │ Meeting с клиентом            │ 25.08.2025 14:00  │
   │ 2 │ Meeting с командой            │ 26.08.2025 10:00  │
   └───┴───────────────────────────────┴───────────────────┘

3. Введите номер события (например: "2").
4. Команда применится к выбранному событию.

─── Различие по вводимым данным:
  • remove / stop_rm / remove_rm → только номер
  • update → номер + новые данные ` + errUpdateFormat + `
`

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

func (c *Cmd) notifyError(err error) bool {
	if err != nil {
		c.handlePrint(err.Error())
		logger.LogError(err.Error())
		return false
	}
	return true
}

func (c *Cmd) handleAddCmd(parts []string) {
	if len(parts) < 4 {
		c.handlePrint(errAddFormat)
		logger.LogError(errAddFormat)
		return
	}
	title := parts[1]
	date := parts[2]
	priority := events.Priority(parts[3])
	msg, err := c.calendar.AddEvent(title, date, priority)
	switch {
	case errors.Is(err, validators.ErrEmptyTitle):
		c.handlePrint(errEmptyTitle)
		logger.LogError(errEmptyTitle)
	case errors.Is(err, validators.ErrDateAlreadyPassed):
		c.handlePrint(errPastTimeTravel)
		logger.LogError(errPastTimeTravel)
	default:
		if !c.notifyResult(msg, err) {
			return
		}
	}
}

func (c *Cmd) handleDeleteCmd(title []string) {
	event, err := c.selectEvents(title)
	if !c.notifyError(err) {
		return
	}
	msg, err := c.calendar.DeleteEvent(event.ID)
	if !c.notifyResult(msg, err) {
		return
	}
}

func (c *Cmd) handleEditeCmd(parts []string) {
	event, err := c.selectEvents(parts)
	if !c.notifyError(err) {
		return
	}
	parts, err = c.readAndParseInput(errUpdateFormat)
	if !c.notifyError(err) {
		return
	}
	if len(parts) < 3 {
		c.handlePrint(errUpdateFormat)
		logger.LogError(errUpdateFormat)
		return
	}
	newTitle := parts[0]
	newDate := parts[1]
	newPriority := events.Priority(parts[2])
	msg, err := c.calendar.EditEvent(event.ID, newTitle, newDate, newPriority)
	if !c.notifyResult(msg, err) {
		return
	}
}

func (c *Cmd) handleAddReminderCmd(parts []string) {
	event, err := c.selectEventsByReminder(false, parts)
	if !c.notifyError(err) {
		return
	}
	parts, err = c.readAndParseInput(errReminderFormat)
	if !c.notifyError(err) {
		return
	}
	if len(parts) < 2 {
		c.handlePrint(errReminderFormat)
		logger.LogError(errReminderFormat)
		return
	}
	message := parts[0]
	time := parts[1]
	msg, err := c.calendar.SetEventReminder(event.ID, message, time)
	if !c.notifyResult(msg, err) {
		return
	}
}

func (c *Cmd) handleStopReminderCmd(title []string) {
	event, err := c.selectEventsByReminder(true, title)
	if !c.notifyError(err) {
		return
	}
	msg, err := c.calendar.CancelEventReminder(event.ID)
	if !c.notifyResult(msg, err) {
		return
	}
}

func (c *Cmd) handleDeleteReminderCmd(title []string) {
	event, err := c.selectEventsByReminder(true, title)
	if !c.notifyError(err) {
		return
	}
	msg, err := c.calendar.RemoveEventReminder(event.ID)
	if !c.notifyResult(msg, err) {
		return
	}
}

func (c *Cmd) handleShowEventsCmd() {
	c.handlePrint(eventShowMessage)
	c.handlePrint(c.calendar.ShowEvents())
}

func (c *Cmd) handleShowLogsCmd() {
	fmt.Println(c.logger.showLogs())
}

func (c *Cmd) handleShowHelpCmd() {
	fmt.Print(helpMessage)
}

func (c *Cmd) handlePrint(msg string) {
	fmt.Println(msg)
	c.logger.logMessage(msg)
}

func (c *Cmd) handleExitCmd() {
	c.calendar.Close()
	os.Exit(0)
}

func (c *Cmd) handleDefaultCmd(cmd string) {
	msg := (unknownCommand + cmd)
	c.handlePrint(msg)
	c.handlePrint(deafaultMessage)
	logger.LogError(msg)
}
