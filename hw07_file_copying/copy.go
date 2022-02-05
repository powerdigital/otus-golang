package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFileSize   = errors.New("unsupported file size")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, srcErr := os.Open(fromPath)
	if srcErr != nil {
		return srcErr
	}

	stat, err := srcFile.Stat()
	if err != nil {
		return ErrUnsupportedFileSize
	}

	srcSize := stat.Size()

	if offset > srcSize || srcSize == 0 {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > srcSize {
		limit = srcSize
	}

	destFile, destErr := os.Create(toPath)
	if destErr != nil {
		return destErr
	}

	defer srcFile.Close()
	defer destFile.Close()

	srcFile.Seek(offset, io.SeekStart)
	io.CopyN(destFile, srcFile, limit)

	return nil
}
