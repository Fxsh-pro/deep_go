package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircularQueuePointers(t *testing.T) {
	const queueSize = 3
	queue, err := NewCircularQueueP(queueSize)
	require.NoError(t, err)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.False(t, queue.Full(), "queue with one element must not be full")
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	require.Equal(t, 1, queue.Front())
	require.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))
	assert.True(t, queue.Full(), "queue must be full after adding a replacement element")

	require.Equal(t, 2, queue.Front())
	require.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	require.Equal(t, 3, queue.Front())
	assert.True(t, queue.Pop())
	require.Equal(t, 4, queue.Front())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
}
