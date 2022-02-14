package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	cmd := []string{
		"./go-envdir",
		"testdata/env",
		"/bin/bash",
		"testdata/echo.sh",
		"arg1=1",
		"arg2=2",
	}

	envs, err := ReadDir(envDir)
	require.NoError(t, err)

	returnCode := RunCmd(cmd, envs)
	require.Equal(t, returnCode, successCode)
}

func TestRunCmdWithoutEnvs(t *testing.T) {
	cmd := []string{
		"./go-envdir",
		"testdata/env",
		"/bin/bash",
		"testdata/echo.sh",
		"arg1=1",
		"arg2=2",
	}

	var envs Environment
	returnCode := RunCmd(cmd, envs)
	require.Equal(t, returnCode, successCode)
}

func TestRunCmdError(t *testing.T) {
	var cmd []string
	envs, err := ReadDir(envDir)
	require.NoError(t, err)

	returnCode := RunCmd(cmd, envs)
	require.Equal(t, returnCode, failureCode)
}
