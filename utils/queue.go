package utils

import "sync"

type queueNode struct {
	data interface{}
	next *queueNode
}

type queue struct {
	head  *queueNode
	tail  *queueNode
	count int
	lock  *sync.Mutex
}

func NewQueue() *queue {
	return &queue{lock: &sync.Mutex{}}
}

func (q *queue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.count
}

func (q *queue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	node := &queueNode{data: v}

	if q.tail == nil {
		q.tail = node
		q.head = node
	} else {
		q.tail.next = node
		q.tail = node
	}

	q.count++
}

func (q *queue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return nil
	}

	node := q.head
	q.head = node.next

	if q.head == nil {
		q.tail = nil
	}

	q.count--

	return node.data
}

func (q *queue) Peek() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	node := q.head
	if node == nil || node.data == nil {
		return nil
	}

	return node.data
}
