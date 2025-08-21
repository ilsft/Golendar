package reminder

import (
	"errors"
	"fmt"
	"time"

	validators "github.com/ilsft/app/utils"
)

const (
	invalidNameRemMsg           = "неверное имя"
	alreadySentRemMsg           = "напоминание уже отправлено!"
	sentRemMsg                  = "напоминание: %s"
	passedTimeRemMsg            = "время напоминания уже прошло"
	remAfterDurationMsg         = "напоминание через: %s"
	remTimerAbsentMsg           = "таймер отсутствует"
	remTimerStoppedMsg          = "таймер остановлен для напоминания: %s"
	remTimerExpiredOrStoppedMsg = "таймер уже сработал или был остановлен: %s"
)

type Reminder struct {
	Message  string      `json:"message"`
	At       time.Time   `json:"time"`
	Sent     bool        `json:"sent"`
	timer    *time.Timer `json:"-"`
	notifier Notifier    `json:"-"`
}

type Notifier interface {
	Notify(msg string)
}

func NewReminder(message string, at time.Time, notifier Notifier) (*Reminder, error) {
	err := validators.CheckTitleEmpty(message)
	if err != nil {
		return nil, err
	}
	if !validators.IsValidTitle(message) {
		return nil, errors.New(invalidNameRemMsg)
	}

	return &Reminder{
		Message:  message,
		At:       at,
		Sent:     false,
		notifier: notifier,
	}, nil

}

func (r *Reminder) Send() {
	if r.Sent {
		r.notifier.Notify(alreadySentRemMsg)
		return
	}
	r.Sent = true
	msg := fmt.Sprintf(sentRemMsg, r.Message)
	r.notifier.Notify(msg)
}

func (r *Reminder) Start() string {
	t := time.Now()
	duration := r.At.Sub(t)
	if duration <= 0 {
		return (passedTimeRemMsg)
	}
	r.timer = time.AfterFunc(duration, func() {
		r.Send()
	})
	return fmt.Sprintf(remAfterDurationMsg, duration.String())
}

func (r *Reminder) Stop() string {
	if r == nil || r.timer == nil {
		return (remTimerAbsentMsg)
	}
	if r.timer.Stop() {
		return fmt.Sprintf(remTimerStoppedMsg, r.Message)
	} else {
		return fmt.Sprintf(remTimerExpiredOrStoppedMsg, r.Message)
	}
}
