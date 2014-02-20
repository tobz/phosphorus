package utils

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewQueue(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)
}

func TestEmptyPush(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)

	q.Push(struct{}{})
}

func TestEmptyPop(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)

	v := q.Pop()
	assert.Nil(t, v)
}

func TestNonEmptyPush(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)

	q.Push(struct{}{})
	q.Push(struct{}{})
}

func TestNonEmptyPop(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)

	q.Push(struct{}{})

	v := q.Pop()
	assert.NotNil(t, v)
}

func TestRandomPushPop(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)

	v := q.Pop()
	assert.Nil(t, v)

	q.Push(struct{}{})
	q.Push(struct{}{})

	v = q.Pop()
	assert.NotNil(t, v)

	q.Push(struct{}{})

	v = q.Pop()
	assert.NotNil(t, v)
	v = q.Pop()
	assert.NotNil(t, v)
	v = q.Pop()
	assert.Nil(t, v)
}

func TestPeek(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)

	q.Push(struct{}{})

	v := q.Pop()
	assert.NotNil(t, v)
	v = q.Pop()
	assert.Nil(t, v)

	q.Push(struct{}{})

	v = q.Peek()
	assert.NotNil(t, v)

	v = q.Pop()
	assert.NotNil(t, v)

	v = q.Peek()
	assert.Nil(t, v)
}

func TestLen(t *testing.T) {
	q := NewQueue()
	assert.NotNil(t, q)

	q.Push(struct{}{})

	l := q.Len()
	assert.Equal(t, 1, l)
	v := q.Pop()
	assert.NotNil(t, v)
	l = q.Len()
	assert.Equal(t, 0, l)

	q.Push(struct{}{})
	q.Push(struct{}{})

	l = q.Len()
	assert.Equal(t, 2, l)

	v = q.Peek()
	assert.NotNil(t, v)

	l = q.Len()
	assert.Equal(t, 2, l)
}
