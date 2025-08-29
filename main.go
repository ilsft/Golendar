package main

import (
	"fmt"

	"github.com/ilsft/Golendar/calendar"
	"github.com/ilsft/Golendar/cmd"
	"github.com/ilsft/Golendar/logger"
	"github.com/ilsft/Golendar/storage"
)

func main() {
	s := storage.NewJsonStorage("calendar.json")
	c := calendar.NewCalendar(s)
	err := c.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	file, err := logger.StartLogger("app.log")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	historyStorage := storage.NewJsonStorage("iohistory.json")
	historyLogger := cmd.NewHistoryLogger(historyStorage)

	cli := cmd.NewCmd(c, historyLogger)
	cli.Run()
}
