package main

import (
	"fmt"

	"github.com/cipepser/go-example/functionalOptions"
)

func main() {
	r := functionalOptions.NewRequest()
	fmt.Println(r) // &{1 30 desc}

	r = functionalOptions.NewRequest(functionalOptions.Page(10))
	fmt.Println(r) // &{10 30 desc}

	r = functionalOptions.NewRequest(
		functionalOptions.Page(10),
		functionalOptions.PerPage(2),
		functionalOptions.Sort("ast"),
	)
	fmt.Println(r) // &{10 2 ast}
}
