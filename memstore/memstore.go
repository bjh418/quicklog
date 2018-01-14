package memstore

import (
	"errors"
	"fmt"
	"time"
	"github.com/devxfactor/quicklog/shared"
	"github.com/devxfactor/quicklog/utils"
)

type Memstore interface {
	Trace(time time.Time, line string)
	Tracef(time time.Time, format string, args ...interface{})
	Debug(time time.Time, line string)
	Debugf(time time.Time, format string, args ...interface{})
	Info(time time.Time, line string)
	Infof(time time.Time, format string, args ...interface{})
	Note(time time.Time, line string)
	Notef(time time.Time, format string, args ...interface{})
	Warn(time time.Time, line string)
	Warnf(time time.Time, format string, args ...interface{})
	Error(time time.Time, line string)
	Errorf(time time.Time, format string, args ...interface{})

	Log(time time.Time, level Level, line string) error
	Logf(time time.Time, level Level, format string, args ...interface{}) error

	Tail(f func (Entry) error) error
}

type Entry interface {
	Time() time.Time
	Level() Level
	Line() string
}

type Level string

const (
	TRACE = iota
	DEBUG
	INFO
	NOTE
	WARN
	ERROR
)
var (
	LEVELS = []Level{ "TRACE", "DEBUG", "INFO", "NOTE", "WARN", "ERROR" }
	levels []string
)
func init() {
	for _, l := range LEVELS {
		levels = append(levels, string(l))
	}
}

type memstore struct {
	list shared.Log
}

type entry struct {
	time time.Time
	level int
	line string
}

func NewMemstore() Memstore {
	list, _ := shared.NewLog(10000)
	m := &memstore{list: list}
	return m
}

func (m *memstore) Trace(time time.Time, line string) {
	m.log(time, TRACE, line)
}
func (m *memstore) Tracef(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	m.Trace(time, line)
}

func (m *memstore) Debug(time time.Time, line string) {
	m.log(time, DEBUG, line)
}
func (m *memstore) Debugf(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	m.Debug(time, line)
}

func (m *memstore) Info(time time.Time, line string) {
	m.log(time, INFO, line)
}
func (m *memstore) Infof(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	m.Info(time, line)
}

func (m *memstore) Note(time time.Time, line string) {
	m.log(time, NOTE, line)
}
func (m *memstore) Notef(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	m.Note(time, line)
}

func (m *memstore) Warn(time time.Time, line string) {
	m.log(time, WARN, line)
}
func (m *memstore) Warnf(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	m.Warn(time, line)
}

func (m *memstore) Error(time time.Time, line string) {
	m.log(time, ERROR, line)
}
func (m *memstore) Errorf(time time.Time, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	m.Error(time, line)
}

func (m *memstore) Logf(time time.Time, level Level, format string, args... interface{}) error {
	line := fmt.Sprintf(format, args...)
	return m.Log(time, level, line)
}

func (m *memstore) Log(time time.Time, level Level, line string) error {
	l, err := levelInt(level)
	if err != nil {
		return err
	}
	m.log(time, l, line)
	return nil
}

func (m *memstore) log(time time.Time, level int, line string) {
	m.list.Add(entry{time: time, level: level, line: line})
}

func (m *memstore) Tail(f func (Entry) error) error {
	return m.list.Tail(func (value interface{}) error {
		return f(value.(Entry))
	})
}

func (e entry) Time() time.Time {
	return e.time
}

func (e entry) Level() Level {
	return Level(LEVELS[e.level])
}

func (e entry) Line() string {
	return e.line
}

func levelInt(level Level) (int, error) {
	index, err := utils.StringIndex(string(level), levels)
	if err != nil {
		return -1, errors.New("invalid level")
	}
	return index, nil
}