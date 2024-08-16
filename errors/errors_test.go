package errors_test

import (
	"fmt"
	"os"

	"github.com/AdguardTeam/golibs/errors"
)

func ExampleAs() {
	if _, err := os.Open("non-existing"); err != nil {
		var pathError *os.PathError

		if errors.As(err, &pathError) {
			fmt.Println("Failed at path:", pathError.Path)
		} else {
			fmt.Println(err)
		}
	}

	// Output:
	//
	// Failed at path: non-existing
}

func ExampleIs() {
	if _, err := os.Open("non-existing"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("file does not exist")
		} else {
			fmt.Println(err)
		}
	}

	// Output:
	//
	// file does not exist
}

func ExampleJoin() {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err := errors.Join(err1, err2)
	fmt.Println(err)
	if errors.Is(err, err1) {
		fmt.Println("err is err1")
	}
	if errors.Is(err, err2) {
		fmt.Println("err is err2")
	}

	// Output:
	// err1
	// err2
	// err is err1
	// err is err2
}

func ExampleNew() {
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}

	// Output:
	//
	// emit macho dwarf: elf header corrupted
}
