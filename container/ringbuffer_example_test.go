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
	fmt.Printf("append %d: %#v\n", x, rb.Current())

	rb.Push(y)
	fmt.Printf("append %d: %#v\n", y, rb.Current())

	rb.Push(z)
	fmt.Printf("append %d: %#v\n", z, rb.Current())
	fmt.Printf("current: %#v\n", rb.Current())

	// Output:
	// nil: 0
	// empty: 0
	// append 1: 0
	// append 2: 1
	// append 3: 2
	// current: 2
}
