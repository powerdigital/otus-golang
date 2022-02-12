package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	correctValues := map[string]string{
		"BAR":   "bar",
		"EMPTY": "",
		"FOO":   "   foo\nwith new line",
		"HELLO": "\"hello\"",
		"UNSET": "",
	}

	envs, err := ReadDir(envDir)
	require.NoError(t, err)

	for envName, env := range envs {
		require.Equal(t, correctValues[envName], env.Value)
	}
}

func TestDirNotExists(t *testing.T) {
	_, err := ReadDir("unknown_dir/")
	require.Error(t, err)
}
