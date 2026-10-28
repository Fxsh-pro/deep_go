package main

import (
	"reflect"
	"runtime"
	"sync"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	*cowBufferState
}

type cowBufferState struct {
	data []byte
	refs *int
	mu   *sync.Mutex
}

func NewCOWBuffer(data []byte) COWBuffer {
	var refs = 1
	state := &cowBufferState{
		data: data,
		refs: &refs,
		mu:   new(sync.Mutex),
	}
	runtime.SetFinalizer(state, finalizeCOWBuffer)
	return COWBuffer{cowBufferState: state}
}

func finalizeCOWBuffer(state *cowBufferState) {
	buffer := COWBuffer{cowBufferState: state}
	buffer.Close()
}

func (b *COWBuffer) Clone() COWBuffer {
	if b.cowBufferState == nil {
		return COWBuffer{}
	}
	defer runtime.KeepAlive(b.cowBufferState)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.refs == nil {
		return COWBuffer{}
	}
	*b.refs++
	state := &cowBufferState{
		data: b.data,
		refs: b.refs,
		mu:   b.mu,
	}
	runtime.SetFinalizer(state, finalizeCOWBuffer)
	return COWBuffer{cowBufferState: state}
}

func (b *COWBuffer) Close() {
	if b.cowBufferState == nil {
		return
	}
	defer runtime.KeepAlive(b.cowBufferState)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.refs == nil {
		return
	}
	*b.refs--
	b.data = nil
	b.refs = nil
	runtime.SetFinalizer(b.cowBufferState, nil)
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if b.cowBufferState == nil {
		return false
	}
	defer runtime.KeepAlive(b.cowBufferState)
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.refs == nil || index < 0 || index >= len(b.data) {
		return false
	}
	if *b.refs > 1 {
		newData := make([]byte, len(b.data))
		copy(newData, b.data)
		*b.refs--
		b.data = newData
		refs := 1
		b.refs = &refs
	}
	b.data[index] = value
	return true
}

func (b *COWBuffer) String() string {
	if b.cowBufferState == nil {
		return ""
	}
	defer runtime.KeepAlive(b.cowBufferState)
	b.mu.Lock()
	defer b.mu.Unlock()
	str := *(*string)(unsafe.Pointer(&b.data))
	return str
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}
