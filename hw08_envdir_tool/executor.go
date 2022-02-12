package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 4 {
		return failureCode
	}

	executableCmd := cmd[3]
	args := strings.Join(cmd[4:], " ")
	command := exec.Command(executableCmd, args)

	for envName, x := range env {
		os.Setenv(envName, x.Value)
	}

	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		fmt.Println(err.Error())
		return failureCode
	}

	return successCode
}
