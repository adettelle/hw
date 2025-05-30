package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse" //nolint:depguard
)

func main() {
	s := "Hello, OTUS!"
	fmt.Println(reverse.String(s))
}
