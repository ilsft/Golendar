package validators

import (
	"fmt"
	"testing"
)

func TestValidTitle(t *testing.T) {
	title := "тест-2"
	success := IsValidTitle(title)
	if !success {
		t.Errorf("Неверное имя: %s", title)
	} else {
		fmt.Printf("Корректное имя: %s\n", title)
	}

	title2 := "тесW"
	success2 := IsValidTitle(title2)
	if !success2 {
		t.Errorf("Неверное имя: %s", title2)
	} else {
		fmt.Printf("Корректное имя: %s\n", title2)
	}
}

func TestValidDate(t *testing.T) {

	title := "2020-12-12"
	time, err := ValidateDate(title)
	if err != nil {
		t.Errorf("Ошибка %v", err)
	} else {
		fmt.Printf("Корректная дата: %v\n", time)
	}

	title2 := "2025-12-12"
	time2, err2 := ValidateDate(title2)
	if err2 != nil {
		t.Errorf("Ошибка %v", err2)
	} else {
		fmt.Printf("Корректная дата: %v\n", time2)
	}
}
