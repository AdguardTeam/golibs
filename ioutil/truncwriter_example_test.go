package ioutil_test

import (
	"bytes"
	"fmt"

	"github.com/AdguardTeam/golibs/ioutil"
)

func ExampleTruncatedWriter() {
	data := []byte("hello")

	buf := &bytes.Buffer{}
	fmt.Println(ioutil.NewTruncatedWriter(buf, 1).Write(data))
	fmt.Println(buf)

	buf = &bytes.Buffer{}
	fmt.Println(ioutil.NewTruncatedWriter(buf, 2).Write(data))
	fmt.Println(buf)

	buf = &bytes.Buffer{}
	fmt.Println(ioutil.NewTruncatedWriter(buf, 10).Write(data))
	fmt.Println(buf)

	// Output:
	//
	// 5 <nil>
	// h
	// 5 <nil>
	// he
	// 5 <nil>
	// hello
}
