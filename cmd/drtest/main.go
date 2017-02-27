package main

import (
	"fmt"

	"github.com/jakebailey/dr"
)

func main() {
	for _, s := range []string{
		"asdfg",
		"aaa+bbb",
		"!(a)b*(cd)*e+f",
		`\+\++\*(\!\\)`,
	} {
		testR := dr.MustParse(s)
		fmt.Printf("%v => %v\n", s, testR)
	}

	fmt.Println()

	r := dr.Concat{
		L: dr.Char{R: 'a'},
		R: dr.Concat{
			L: dr.Char{R: 'b'},
			R: dr.Kleene{
				R: dr.Char{R: 'c'},
			},
		},
	}

	fmt.Println("matching against", r)

	for _, s := range []string{
		"",
		"a",
		"ab",
		"abc",
		"abccccc",
	} {
		fmt.Printf("%v: %v\n", s, dr.Match(r, s))
	}
}
