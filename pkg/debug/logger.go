package debug

import (
	"fmt"
	"sync"
	"time"
)

type LogLevel int

const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarn
	LogError
)

func (l LogLevel) String() string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Category  string
	Message   string
	Context   map[string]interface{}
}

func (e LogEntry) String() string {
	timeStr := e.Timestamp.Format("15:04:05.000")
	if len(e.Context) == 0 {
		return fmt.Sprintf("[%s] %s [%s] %s", timeStr, e.Level, e.Category, e.Message)
	}
	return fmt.Sprintf("[%s] %s [%s] %s %+v", timeStr, e.Level, e.Category, e.Message, e.Context)
}

type Logger struct {
	mu      sync.RWMutex
	entries []LogEntry
	maxSize int
	enabled bool
}

var globalLogger = &Logger{
	maxSize: 1000,
	enabled: true,
	entries: make([]LogEntry, 0, 1000),
}

func SetEnabled(enabled bool) {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.enabled = enabled
}

func SetMaxSize(size int) {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.maxSize = size
	if len(globalLogger.entries) > size {
		copy(globalLogger.entries, globalLogger.entries[len(globalLogger.entries)-size:])
		globalLogger.entries = globalLogger.entries[:size]
	}
}

func (l *Logger) log(level LogLevel, category, message string, context map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.enabled {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Category:  category,
		Message:   message,
		Context:   context,
	}

	if len(l.entries) >= l.maxSize {
		copy(l.entries, l.entries[1:])
		l.entries[l.maxSize-1] = entry
	} else {
		l.entries = append(l.entries, entry)
	}
}

func (l *Logger) GetEntries() []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entries := make([]LogEntry, len(l.entries))
	copy(entries, l.entries)
	return entries
}

func (l *Logger) GetEntriesFiltered(level LogLevel, category string) []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var filtered []LogEntry
	for _, entry := range l.entries {
		if entry.Level >= level {
			if category == "" || entry.Category == category {
				filtered = append(filtered, entry)
			}
		}
	}
	return filtered
}

func (l *Logger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = l.entries[:0]
}

func Debug(category, message string) {
	globalLogger.log(LogDebug, category, message, nil)
}

func Debugf(category, format string, args ...interface{}) {
	globalLogger.log(LogDebug, category, fmt.Sprintf(format, args...), nil)
}

func DebugWith(category, message string, context map[string]interface{}) {
	globalLogger.log(LogDebug, category, message, context)
}

func Info(category, message string) {
	globalLogger.log(LogInfo, category, message, nil)
}

func Infof(category, format string, args ...interface{}) {
	globalLogger.log(LogInfo, category, fmt.Sprintf(format, args...), nil)
}

func InfoWith(category, message string, context map[string]interface{}) {
	globalLogger.log(LogInfo, category, message, context)
}

func Warn(category, message string) {
	globalLogger.log(LogWarn, category, message, nil)
}

func Warnf(category, format string, args ...interface{}) {
	globalLogger.log(LogWarn, category, fmt.Sprintf(format, args...), nil)
}

func WarnWith(category, message string, context map[string]interface{}) {
	globalLogger.log(LogWarn, category, message, context)
}

func Error(category, message string) {
	globalLogger.log(LogError, category, message, nil)
}

func Errorf(category, format string, args ...interface{}) {
	globalLogger.log(LogError, category, fmt.Sprintf(format, args...), nil)
}

func ErrorWith(category, message string, context map[string]interface{}) {
	globalLogger.log(LogError, category, message, context)
}

func GetEntries() []LogEntry {
	return globalLogger.GetEntries()
}

func GetEntriesFiltered(level LogLevel, category string) []LogEntry {
	return globalLogger.GetEntriesFiltered(level, category)
}

func Clear() {
	globalLogger.Clear()
}

func GetCategories() []string {
	globalLogger.mu.RLock()
	defer globalLogger.mu.RUnlock()

	categorySet := make(map[string]bool)
	for _, entry := range globalLogger.entries {
		categorySet[entry.Category] = true
	}

	categories := make([]string, 0, len(categorySet))
	for category := range categorySet {
		categories = append(categories, category)
	}
	return categories
}
