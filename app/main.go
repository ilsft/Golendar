package main

import (
	"fmt"

	"github.com/ilsft/app/calendar"
)

func main() {
	// ошибки при добавлении
	calendar.AddEvent("Fishing15Ё!", "2025-10-10 12:00", "High")
	calendar.AddEvent("Fishing", "2025-10-10 122:00", "Low")
	// корректное удаление
	calendar.AddEvent("Картошка", "2025-10-11 08:00", "High")
	calendar.AddEvent("Свадьба", "2025-10-12 09:00", "Low")
	calendar.AddEvent("Бадминтон", "2025-10-10 12:00", "Low")

	// ошибка при удалении
	calendar.DeleteEvent("Бадмифынтон")
	// корректное удаление
	calendar.DeleteEvent("Бадминтон")

	//ошибки при изменении
	calendar.EditEvent("Свдьба", "2025-09-12", "High")
	calendar.EditEvent("Свадьба", "2025-0912", "High")
	calendar.EditEvent("Свадьба", "2025-09-12", "Higывh")
	// корректное изменение
	calendar.EditEvent("Свадьба", "2025-09-12", "High")
	calendar.ShowEvents()
	fmt.Scanln()
}
