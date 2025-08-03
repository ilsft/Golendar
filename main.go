package main

import (
	"fmt"

	"github.com/ilsft/app/calendar"
)

const ()

func main() {

	calendar.AddEvent("Свадьба", "2025-10-12 09:00", "Low")
	calendar.AddEvent("Бадминтон", "2025-10-10 12:00", "Low")

	event1, err := calendar.AddEvent("Картошка", "2025-10-11 08:00", "High")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(event1)

	event2, err2 := calendar.EditTitleEvent("Свадьба", "Почесать пузо")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	fmt.Println(event2)

	event3, err3 := calendar.EditDateEvent("Картошка", "2026-10-11 08:00")
	if err3 != nil {
		fmt.Println(err3)
		return
	}
	fmt.Println(event3)

	event4, err4 := calendar.EditPriorityEvent("Картошка", "Low")
	if err4 != nil {
		fmt.Println(err4)
		return
	}
	fmt.Println(event4)

	event5, err5 := calendar.DeleteEvent("Бадминтон")
	if err5 != nil {
		fmt.Println(err5)
		return
	}
	fmt.Println(event5)

	calendar.ShowEvents()
	fmt.Scanln()
}
