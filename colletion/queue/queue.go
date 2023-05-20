package queue

import (
	"errors"
)

// Queue 队列
type Queue struct {
	items []interface{}
	size  int
}

func NewQueue() *Queue {
	return &Queue{
		items: make([]interface{}, 0),
		size:  0,
	}
}

func (q *Queue) isEmpty() bool {
	return q.size == 0
}

func (q *Queue) getSize() int {
	return q.size
}

func (q *Queue) push(val interface{}) {
	q.items = append(q.items, val)
	q.size++
}

func (q *Queue) pop() (data interface{}, err error) {
	if q.isEmpty() {
		err = errors.New("empty queue")
	}
	data = q.items[0]
	q.items = q.items[1:]
	q.size--
	return data, err
}

func (q *Queue) PeekNext() (data interface{}, err error) {
	if q.isEmpty() {
		err = errors.New("empty queue")
	}
	return q.items[0], err
}
