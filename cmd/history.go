package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ilsft/app/storage"
)

type LogEntry struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type Logger struct {
	Logs    []LogEntry
	Storage storage.Store
}

var mu sync.Mutex

var logsStorage = storage.NewJsonStorage("iohistory.json")
var l = NewLogs(logsStorage)

func NewLogs(s storage.Store) *Logger {
	return &Logger{
		Logs:    make([]LogEntry, 0),
		Storage: s,
	}
}

func (l *Logger) logMessage(message string) {
	mu.Lock()
	defer mu.Unlock()
	t := time.Now()
	entry := LogEntry{
		Time:    t,
		Message: message,
	}
	l.Logs = append(l.Logs, entry)
}

func (l *Logger) loadLogs() error {
	data, err := l.Storage.Load()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		if l.Logs == nil {
			l.Logs = make([]LogEntry, 0)
		}
		return nil
	}
	err = json.Unmarshal(data, l)
	if err != nil {
		return (err)
	}
	return nil
}

func (l *Logger) saveLogs() error {
	data, err := json.Marshal(l)
	if err != nil {
		return (err)
	}
	err = l.Storage.Save(data)
	if err != nil {
		return (err)
	}
	return nil
}

func (l *Logger) showLogs() string {
	if len(l.Logs) == 0 {
		return ""
	}
	var logs []string
	for _, entry := range l.Logs {
		log := (fmt.Sprintf("%s - %s", entry.Time.Format("2006-01-02 15:04:05"), entry.Message))
		logs = append(logs, log)
	}
	return strings.Join(logs, "\n")
}
