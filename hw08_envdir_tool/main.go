package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("incorrect number of arguments")
		os.Exit(-1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("unable to read path to environment, %v", err)
		os.Exit(-1)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
