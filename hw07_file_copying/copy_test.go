package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

const srcFolder = "testdata"
const destFolder = "/tmp"

func TestCopySuccess(t *testing.T) {
	dirs, err := ioutil.ReadDir(srcFolder)
	require.NoError(t, err)

	for _, f := range dirs {
		src := fmt.Sprintf("%s/%s", srcFolder, f.Name())
		dest := fmt.Sprintf("%s/%s", destFolder, f.Name())
		err := Copy(src, dest, 0, 0)
		require.Nil(t, err)
		require.FileExists(t, dest)
	}
}

func TestOffsetExceedsFileSize(t *testing.T) {
	filename := "out_offset0_limit10.txt"
	src := fmt.Sprintf("%s/%s", srcFolder, filename)
	dest := fmt.Sprintf("%s/%s", destFolder, filename)
	err := Copy(src, dest, 100, 0)
	require.Error(t, err)
	require.Equal(t, err, ErrOffsetExceedsFileSize)
}

func TestLimitExceedsFileSize(t *testing.T) {
	filename := "out_offset0_limit10.txt"
	src := fmt.Sprintf("%s/%s", srcFolder, filename)
	dest := fmt.Sprintf("%s/%s", destFolder, filename)
	err := Copy(src, dest, 0, 100)
	require.NoError(t, err)
}
