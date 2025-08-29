package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ilsft/Golendar/storage"
)

type HistoryEntry struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type HistoryLogger struct {
	Logs    []HistoryEntry
	Storage storage.Store
}

var mu sync.Mutex

func NewHistoryLogger(s storage.Store) *HistoryLogger {
	return &HistoryLogger{
		Logs:    make([]HistoryEntry, 0),
		Storage: s,
	}
}

func (hl *HistoryLogger) logMessage(message string) {
	mu.Lock()
	defer mu.Unlock()
	entry := HistoryEntry{
		Time:    time.Now(),
		Message: message,
	}
	hl.Logs = append(hl.Logs, entry)
}

func (hl *HistoryLogger) loadLogs() error {
	data, err := hl.Storage.Load()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		if hl.Logs == nil {
			hl.Logs = make([]HistoryEntry, 0)
		}
		return nil
	}
	err = json.Unmarshal(data, hl)
	if err != nil {
		return (err)
	}
	return nil
}

func (hl *HistoryLogger) saveLogs() error {
	data, err := json.Marshal(hl)
	if err != nil {
		return (err)
	}
	err = hl.Storage.Save(data)
	if err != nil {
		return (err)
	}
	return nil
}

func (hl *HistoryLogger) showLogs() string {
	if len(hl.Logs) == 0 {
		return ""
	}
	var logs []string
	for _, entry := range hl.Logs {
		log := (fmt.Sprintf("%s - %s", entry.Time.Format(patternTime), entry.Message))
		logs = append(logs, log)
	}
	return strings.Join(logs, "\n")
}
