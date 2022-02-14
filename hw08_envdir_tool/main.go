package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Too few arguments, at least 4 required")
		return
	}

	envDir := args[1]

	envs, err := ReadDir(envDir)
	if err != nil {
		return
	}

	RunCmd(args, envs)
}
