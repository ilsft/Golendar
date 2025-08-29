package validators

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

const (
	validPattern = "^[a-zA-Z\\p{Cyrillic}0-9 ]{3,50}$"
	timePattern  = "Mon 2006/01/02 - 15:04"
)

var (
	ErrEmptyTitle        = errors.New("пустая строка содержит только пробелы")
	ErrDateAlreadyPassed = errors.New("указанная дата уже прошла")
)

func FormatDateEvent(date time.Time) string {
	return date.Format(timePattern)
}

func IsValidTitle(title string) bool {
	matched, err := regexp.MatchString(validPattern, title)
	if err != nil {
		return false
	}
	return matched
}
func CheckTitleEmpty(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}
	return nil
}

func ValidateDate(dateStr string) (time.Time, error) {
	t, err := dateparse.ParseAny(dateStr)
	if err != nil {
		return time.Time{}, err
	}

	localTime := time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		time.Local,
	)

	if !localTime.After(time.Now()) {
		return time.Time{}, ErrDateAlreadyPassed
	}

	return localTime,
		nil
}
