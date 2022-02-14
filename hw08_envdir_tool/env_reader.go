package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)
	for _, f := range files {
		filepath := fmt.Sprintf("%s/%s", dir, f.Name())
		info, err := os.Stat(filepath)
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			result[f.Name()] = EnvValue{"", true}
		}

		file, err := os.Open(filepath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			envVal := strings.TrimRight(scanner.Text(), "\t\n ")
			needRemove := len(envVal) == 0
			source := []byte(envVal)
			envVal = string(bytes.Replace(source, []byte{0x00}, []byte{'\n'}, 1))

			result[f.Name()] = EnvValue{envVal, needRemove}

			break
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return result, nil
}
