package main

import (
	"fmt"

	"github.com/ilsft/Golendar/app/calendar"
	"github.com/ilsft/Golendar/app/cmd"
	"github.com/ilsft/Golendar/app/logger"
	"github.com/ilsft/Golendar/app/storage"
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

	cli := cmd.NewCmd(c)
	cli.Run()
}
