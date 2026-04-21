package ioutil_test

import (
	"fmt"
	"slices"
	"sync"

	"github.com/AdguardTeam/golibs/errors"
	"github.com/AdguardTeam/golibs/ioutil"
)

func ExampleLockedBuffer() {
	const goroutineCount = 5

	b := ioutil.NewLockedBuffer(goroutineCount)

	var wg sync.WaitGroup

	for i := range goroutineCount {
		wg.Go(func() {
			_, err := b.Write([]byte{byte(i)})
			errors.Check(err)
		})
	}

	wg.Wait()

	readBuffer := make([]byte, goroutineCount)

	for i := range goroutineCount {
		wg.Go(func() {
			_, err := b.Read(readBuffer[i : i+1])
			errors.Check(err)
		})
	}

	wg.Wait()

	slices.Sort(readBuffer)

	fmt.Printf("%x\n", readBuffer)

	// Output:
	//
	// 0001020304
}
