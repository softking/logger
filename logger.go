package logger

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var LogBufferLength = 4096

type Level int

const (
	DEBUG    = 0 //32
	TRACE    = 1 //16
	INFO     = 2 //8
	WARNING  = 3 //4
	ERROR    = 4 //2
	CRITICAL = 5 //1
)

// Logging level strings
var (
	levelStrings = [...]string{"DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}

	// levelMapping 日志等级表
	levelMapping = map[string]int{
		"DEBUG":    32,
		"TRACE":    16,
		"INFO":     8,
		"WARNING":  4,
		"ERROR":    2,
		"CRITICAL": 1,
	}
)

// LogRecord 日志结构
type LogRecord struct {
	Level   Level     // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Message string    // The log message
}

// LogWriter writer 接口
type LogWriter interface {
	LogWrite(rec *LogRecord)
	Close()
}

// Filter filter
type Filter struct {
	Level    Level
	LevelStr string
	LogWriter
}

type Logger map[string]*Filter

func (log Logger) Close() {
	for name, filt := range log {
		filt.Close()
		delete(log, name)
	}
}

func (log Logger) AddFilter(name string, lvl Level, writer LogWriter) Logger {
	levelStr := strconv.FormatInt(int64(lvl), 2)
	for len(levelStr) < len(levelStrings) {
		levelStr = "0" + levelStr
	}
	log[name] = &Filter{Level: lvl, LevelStr: levelStr, LogWriter: writer}
	return log
}

func (log Logger) intLogf(lvl Level, format string, args ...interface{}) {

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Message: msg,
	}

	for _, filt := range log {
		if filt.LevelStr[lvl] == 48 {
			continue
		}
		filt.LogWrite(rec)
	}

}

func (log Logger) Debug(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(lvl, first, args...)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func (log Logger) Trace(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(lvl, first, args...)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func (log Logger) Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		log.intLogf(lvl, first, args...)
	default:
		// Build a format string so that it will be similar to Sprint
		log.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func (log Logger) Warn(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(lvl, msg)
	return errors.New(msg)
}

func (log Logger) Error(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(lvl, msg)
	return errors.New(msg)
}

func (log Logger) Critical(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	var msg string
	switch first := arg0.(type) {
	case string:
		msg = fmt.Sprintf(first, args...)
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	log.intLogf(lvl, msg)
	return errors.New(msg)
}

// BuildLevel 根据字符串转化为level
func BuildLevel(lvl string) Level {
	lvlStr := strings.Split(lvl, "|")
	num := 0
	for _, l := range lvlStr {
		num += levelMapping[l]
	}
	return Level(num)
}