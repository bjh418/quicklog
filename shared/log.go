package shared

import (
	"sync"
	"errors"
)

type Log interface {
	Add(value interface{})
	Tail(f func(value interface{}) error) error
}

type log struct {
	maxLen int
	len    int
	head   *node
	last   *node
	cond   *sync.Cond
}

type node struct {
	next *node
	value interface{}
}

func NewLog(maxLen int) (Log, error) {
	if maxLen <= 0 {
		return nil, errors.New("maxLen must be positive")
	}
	return &log{maxLen: maxLen, len: 0, head: nil, last: nil, cond: sync.NewCond(&sync.Mutex{})}, nil
}

func (l *log) Add(value interface{}) {
	l.cond.L.Lock()
	last := &node{nil, value}
	if l.head == nil {
		l.head = last
		l.last = last
	} else {
		l.last.next = last
		l.last = last
	}
	if l.len == l.maxLen {
		l.head = l.head.next
	} else {
		l.len += 1
	}
	l.cond.L.Unlock()
	l.cond.Broadcast()
}

func (l *log) Tail(f func(value interface{}) error) error {
	cond := l.cond
	cond.L.Lock()
	for l.head == nil {
		cond.Wait()
	}
	curr := l.head
	cond.L.Unlock()

	for {
		err := f(curr.value)
		if err != nil {
			return err
		}

		if curr.next == nil {
			cond.L.Lock()
			for curr.next == nil {
				cond.Wait()
			}
			cond.L.Unlock()
		}
		curr = curr.next
	}
}