package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %s [input file]\n", args[0])
		os.Exit(-1)
	}
	module := VMExecutionModuleInit()
	file, err := os.Open(args[1])
	if err != nil {
		fmt.Printf("Failed to open the file: %s\n", args[1])
		os.Exit(-1)
	}
	module.BindToFile(file)
	module.Execute()
	file.Close()
}
