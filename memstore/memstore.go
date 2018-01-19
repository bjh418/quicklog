package memstore

import (
	"github.com/devxfactor/quicklog/shared"
)

type Memstore interface {
	Log(line string)
	Tail(f func (string) error) error
}

type memstore struct {
	list shared.Log
}

func NewMemstore() Memstore {
	list, _ := shared.NewLog(10000)
	m := &memstore{list: list}
	return m
}

func (m *memstore) Log(line string) {
	m.list.Add(line)
}

func (m *memstore) Tail(f func (string) error) error {
	return m.list.Tail(func (value interface{}) error {
		return f(value.(string))
	})
}
