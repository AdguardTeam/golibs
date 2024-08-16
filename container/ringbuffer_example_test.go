package container_test

import (
	"fmt"

	"github.com/AdguardTeam/golibs/container"
)

func ExampleRingBuffer_Current() {
	const size, x, y, z = 2, 1, 2, 3

	var rb *container.RingBuffer[int]
	fmt.Printf("nil: %#v\n", rb.Current())

	rb = container.NewRingBuffer[int](size)
	fmt.Printf("empty: %#v\n", rb.Current())

	rb.Push(x)
	fmt.Printf("push %d: %#v\n", x, rb.Current())

	rb.Push(y)
	fmt.Printf("push %d: %#v\n", y, rb.Current())

	rb.Push(z)
	fmt.Printf("push %d: %#v\n", z, rb.Current())
	fmt.Printf("current: %#v\n", rb.Current())

	// Output:
	// nil: 0
	// empty: 0
	// push 1: 0
	// push 2: 1
	// push 3: 2
	// current: 2
}
