package main

import (
	"fmt"
	"testing"
)

func BenchmarkCircularQueue(b *testing.B) {
	for _, size := range []int{1, 5, 100, 100_000, 1_000_000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.Run("array", func(b *testing.B) {
				queue, err := NewCircularQueue(size)
				if err != nil {
					b.Fatal(err)
				}
				validateQueueForBenchmark(b, &queue, size)
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					queue, err := NewCircularQueue(size)
					if err != nil {
						b.Fatal(err)
					}
					for value := 0; value < size; value++ {
						if !queue.Push(value) {
							b.Fatal("Push failed while filling the queue")
						}
					}
					for count := 0; count < size; count++ {
						if !queue.Pop() {
							b.Fatal("Pop failed while draining the queue")
						}
					}
				}
			})

			b.Run("pointers", func(b *testing.B) {
				queue, err := NewCircularQueueP(size)
				if err != nil {
					b.Fatal(err)
				}
				validateQueueForBenchmark(b, &queue, size)
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					queue, err := NewCircularQueueP(size)
					if err != nil {
						b.Fatal(err)
					}
					for value := 0; value < size; value++ {
						if !queue.Push(value) {
							b.Fatal("Push failed while filling the queue")
						}
					}
					for count := 0; count < size; count++ {
						if !queue.Pop() {
							b.Fatal("Pop failed while draining the queue")
						}
					}
				}
			})
		})
	}
}

func validateQueueForBenchmark(b *testing.B, queue interface {
	Push(int) bool
	Pop() bool
	Front() int
	Back() int
	Empty() bool
	Full() bool
}, size int) {
	b.Helper()
	defer func() {
		if err := recover(); err != nil {
			b.Fatalf("queue validation before timing panicked: %v", err)
		}
	}()

	if !queue.Empty() || queue.Full() || queue.Front() != -1 || queue.Back() != -1 {
		b.Fatal("queue validation: new queue must be empty")
	}
	for value := 0; value < size; value++ {
		if !queue.Push(value) {
			b.Fatalf("queue validation: Push(%d) failed", value)
		}
	}
	if queue.Empty() || !queue.Full() || queue.Back() != size-1 {
		b.Fatal("queue validation: incorrect state after filling the queue")
	}
	for want := 0; want < size; want++ {
		if got := queue.Front(); got != want {
			b.Fatalf("queue validation: Front() = %d, want %d", got, want)
		}
		if !queue.Pop() {
			b.Fatalf("queue validation: Pop failed at element %d", want)
		}
	}
	if !queue.Empty() || queue.Full() || queue.Front() != -1 || queue.Back() != -1 {
		b.Fatal("queue validation: incorrect state after draining the queue")
	}
}
